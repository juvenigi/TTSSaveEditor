package resource_map

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"tts-cache-manager-cli/httpfetcher"
	"tts-cache-manager-cli/tabletop_save"
	"tts-cache-manager-cli/util"

	"gopkg.in/yaml.v3"
)

type yamlPackData struct {
	Sha3Archive map[string]string `yaml:"sha-3-map"`
	UrlArchive  map[string]string `yaml:"url-map"`
}

type PackData struct {
	packDataDir   string
	saveTemplates []string
	sha3Archive   map[string]string
	urlArchive    map[string]string
}

func NewPackData(resourceMapFile string) (PackData, error) {
	var result PackData
	var yamlPd yamlPackData
	result.packDataDir = filepath.Dir(resourceMapFile)
	result.sha3Archive = make(map[string]string)
	result.urlArchive = make(map[string]string)

	src, err := os.Open(resourceMapFile)
	if err != nil {
		return result, err
	}

	if entries, err := os.ReadDir(result.packDataDir); err != nil {
		return result, err
	} else {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasPrefix(entry.Name(), ".pack.json") {
				result.saveTemplates = append(result.saveTemplates, entry.Name())
			}
		}
	}

	// todo: could one use a better solution here? (invalid yaml gets skipped without the user knowing)
	if err = yaml.NewDecoder(src).Decode(&yamlPd); err == nil {
		if yamlPd.UrlArchive != nil {
			result.urlArchive = yamlPd.UrlArchive
		} else {
			result.sha3Archive = make(map[string]string)
		}
		if yamlPd.Sha3Archive != nil {
			result.sha3Archive = yamlPd.Sha3Archive
		} else {
			result.sha3Archive = make(map[string]string)
		}

		err = result.cleanup()
	} else {
		// todo: consider a more elegant solution
		//if err := result.FlushToDisk(); err != nil {
		//	return result, err
		//}
	}
	return result, nil
}

func (rm *PackData) FlushToDisk() error {
	file, err := os.Open(rm.packDataDir)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewEncoder(file).Encode(yamlPackData{
		Sha3Archive: rm.sha3Archive,
		UrlArchive:  rm.urlArchive,
	}); err != nil {
		return err
	}

	return file.Sync()
}

// filename, count, extension(may be empty)
var duplicateFilenamePattern = regexp.MustCompile("(.*)-(\\d+)(\\.?.*)$")

func (rm *PackData) AddLocalFile(hashMu *sync.Mutex, res *tabletop_save.GameResource, filenames *RwFilenameSlice) error {
	var location = strings.TrimPrefix(res.ResourceUrl, "file:///")
	blob, sum := util.GetFileAndChecksum(location)
	if sum == nil {
		return errors.New("file not found")
	}

	hashMu.Lock()
	packDataFile, exists := rm.sha3Archive[hex.EncodeToString(sum)]
	hashMu.Unlock()
	if exists {
		res.ResourceUrl = "pack:///" + packDataFile
		res.Status = tabletop_save.Packed
		return nil
	}

	var destFilename = tabletop_save.GetCacheFilename(filepath.Base(res.ResourceUrl))
	return writeResource(filenames, destFilename, rm, blob, res)
}

// avoid repeatedly calling os.ReadDir every time we add a new file.
type RwFilenameSlice struct {
	mu      *sync.RWMutex
	entries map[string]struct{}
}

func NewRwFilenameSlice(packDataDir string) *RwFilenameSlice {
	fileEntries, _ := os.ReadDir(packDataDir)
	filenames := make(map[string]struct{})
	for _, file := range fileEntries {
		if file.IsDir() {
			continue
		}
		filenames[file.Name()] = struct{}{}
	}

	return &RwFilenameSlice{
		mu:      new(sync.RWMutex),
		entries: filenames,
	}
}

// consider an alternative: instead of writing to pack data immediately, aggregate results via a channel (use seen sync.Map for dedupe)
func (rm *PackData) AddRemoteFile(ctx context.Context, hashMu *sync.Mutex, res *tabletop_save.GameResource, filenames *RwFilenameSlice) error {
	urlRes, err := httpfetcher.GetResourceAsync(ctx, res)
	if err != nil {
		res.Status = tabletop_save.Failed
		return err
	}

	hashMu.Lock()
	packFileName, ok := rm.sha3Archive[hex.EncodeToString(urlRes.Checksum)]
	hashMu.Unlock()
	if ok {
		if file, err := os.ReadFile(filepath.Join(rm.packDataDir, packFileName)); err == nil {
			res.Status = tabletop_save.Failed // todo: re-check pack data consistency if such error is detected
			return err

		} else if bytes.Compare(urlRes.Data, file) != 0 {
			res.Status = tabletop_save.Failed
			return errors.New("hash collision / file inconsistency detected (you are a unicorn)")
		}

		res.ResourceUrl = "pack://" + packFileName
		res.Status = tabletop_save.Packed
		return nil
	}
	// first we check the entries ahead of time, then we rely on the singleton nature of the syscalls
	escapedName := tabletop_save.GetCacheFilename(urlRes.Url)
	return writeResource(filenames, escapedName, rm, urlRes.Data, res)
}

func writeResource(filenames *RwFilenameSlice, destinationName string, rm *PackData, blob []byte, res *tabletop_save.GameResource) error {
	filenames.mu.RLock()
	attempt := 0
	attemptedName := fmt.Sprintf("%s-%d", destinationName, attempt)
	// try without attempt counter first, then iterate over attempts
	if _, present := filenames.entries[destinationName]; present {
		for {
			// loop until a unique filename is found
			if _, present = filenames.entries[attemptedName]; !present {
				break
			}
			attempt++
		}
	}
	filenames.mu.RUnlock()

	written, failedAttempts, criticalErr := attemptToWriteExtLess(rm.packDataDir, destinationName, attempt, blob)
	if criticalErr != nil {
		res.Status = tabletop_save.Failed
		return criticalErr
	}
	filenames.mu.Lock()
	filenames.entries[written] = struct{}{}
	for _, fAttempt := range failedAttempts {
		filenames.entries[fAttempt] = struct{}{}
	}
	filenames.mu.Unlock()

	res.ResourceUrl = "pack://" + written
	res.Status = tabletop_save.Packed

	return flushAndUpdate(res, rm, blob, filenames)
}

func flushAndUpdate(res *tabletop_save.GameResource, rm *PackData, file []byte, filenames *RwFilenameSlice) error {
	escapedName := tabletop_save.GetCacheFilename(res.ResourceUrl)
	written, failedAttempts, criticalErr := attemptToWriteExtLess(rm.packDataDir, escapedName, 0, file)
	if criticalErr != nil {
		res.Status = tabletop_save.Failed
		return criticalErr
	}
	filenames.mu.Lock()
	filenames.entries[written] = struct{}{}
	for _, fAttempt := range failedAttempts {
		filenames.entries[fAttempt] = struct{}{}
	}
	filenames.mu.Unlock()

	res.ResourceUrl = "pack://" + written
	res.Status = tabletop_save.Packed

	return nil
}

func attemptToWriteExtLess(dir string, name string, count int, content []byte) (string, []string, error) {
	tries := 0
	var previousAttempts []string
	attemptName := fmt.Sprintf("%s-%d", name, count)
	previousAttempts = append(previousAttempts, attemptName)
	for {
		file, err := os.OpenFile(filepath.Join(dir, attemptName), os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.ModePerm)
		if err == nil {
			_, err := io.Copy(file, bytes.NewReader(content))
			defer file.Close()
			if err != nil {
				return "", previousAttempts, err
			}
			if err := file.Sync(); err != nil {
				return "", previousAttempts, err
			}
			break
		} else if !os.IsExist(err) {
			return "", previousAttempts, err
		}
		if tries++; tries > count {
			break
		} else {
			previousAttempts = append(previousAttempts, attemptName)
			attemptName = fmt.Sprintf("%s-%d", attemptName, count+tries)
		}
	}
	return attemptName, previousAttempts, nil
}

// todo: actually, I would like to already have the filepath of the cached filename here
func (rm *PackData) AddRemoteCachedFile(ctx context.Context, hashMu *sync.Mutex, res *tabletop_save.GameResource, filenames *RwFilenameSlice, scanner *tabletop_save.GameCache) error {
	fpath, ok := scanner.ResourceCached[res.ResourceUrl]
	if !ok {
		return errors.New("cached file not found")
	}

	file, fileChecksum := util.GetFileAndChecksum(fpath)
	if fileChecksum == nil {
		return errors.New("cached file not found")
	}

	hashMu.Lock()
	packDatafile, ok := rm.sha3Archive[hex.EncodeToString(fileChecksum)]
	hashMu.Unlock()

	if ok {
		res.ResourceUrl = "pack://" + packDatafile
		res.Status = tabletop_save.Packed
		return nil
	}

	destinationFilename := tabletop_save.GetCacheFilename(res.ResourceUrl)
	return writeResource(filenames, destinationFilename, rm, file, res)
}

func (rm *PackData) GetPackDataFileFromChecksum(shasum string) (string, bool) {
	res, ok := rm.sha3Archive[shasum]

	return res, ok
}

func (rm *PackData) GetUrlArchiveFromUrl(url string) (string, bool) {
	res, ok := rm.urlArchive[url]

	return res, ok
}

func (rm *PackData) GetSaveTemplates() []string {
	return rm.saveTemplates
}

func (rm *PackData) cleanup() error {
	filenames := NewRwFilenameSlice(rm.packDataDir)

	deletedOnce := false
	deletedOnce = cleanupMap(rm.scanMap(rm.urlArchive, filenames), rm.urlArchive)
	deletedOnce = cleanupMap(rm.scanMap(rm.sha3Archive, filenames), rm.sha3Archive)

	if deletedOnce {
		return rm.FlushToDisk()
	}

	return nil
}

type presentKey struct {
	present bool
	key     string
}

func cleanupMap(mappedUrlExists map[string]presentKey, target map[string]string) bool {
	deletedOnce := false
	for _, key := range mappedUrlExists {
		if key.present {
			delete(target, key.key)
			deletedOnce = true
		}
	}
	return deletedOnce
}

func (rm *PackData) scanMap(urlToFileMap map[string]string, filenames *RwFilenameSlice) map[string]presentKey {
	var mappedFileExists = make(map[string]presentKey)
	for url, file := range urlToFileMap {
		if _, ok := filenames.entries[file]; ok {
			mappedFileExists[file] = presentKey{present: true, key: url}
		} else {
			mappedFileExists[file] = presentKey{present: false, key: url}
		}
	}
	return mappedFileExists
}

func (rm *PackData) AddResourcesFromSave(ctx context.Context, data *tabletop_save.TSSaveFile, gameDir string) error {
	resources := data.GetAllResources()
	scanner, err := tabletop_save.NewCacheScanner(gameDir)
	if err != nil {
		return err
	}

	filenames := NewRwFilenameSlice(rm.packDataDir)
	for _, resource := range resources {
		scanner.SetCacheStatus(resource)
	}

	var errChan = make(chan error)
	// todo: combine ctx with errChan being closed (although we kinda don't care about errors atm)
	//  upd: my plan is to transmit errors as toasts to frontend
	defer close(errChan)

	var hashMu = new(sync.Mutex)
	var wg sync.WaitGroup

	// may contain filenames or urls
	var seen = make(map[string]struct{})
	for idx := range resources {
		var resourceRef = resources[idx]
		if _, already := seen[resourceRef.ResourceUrl]; already {
			continue
		}

		seen[resourceRef.ResourceUrl] = struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch resourceRef.Status {
			case tabletop_save.Remote, tabletop_save.RemoteCached:
				if packUrl, ok := rm.urlArchive[resourceRef.ResourceUrl]; ok {
					resourceRef.ResourceUrl = "pack://" + packUrl
					resourceRef.Status = tabletop_save.Packed
					return
				} else if resourceRef.Status == tabletop_save.Remote {
					_ = rm.AddRemoteFile(ctx, hashMu, resourceRef, filenames)
				} else {
					_ = rm.AddRemoteCachedFile(ctx, hashMu, resourceRef, filenames, &scanner)
				}
			case tabletop_save.Local:
				_ = rm.AddLocalFile(hashMu, resourceRef, filenames)
			default:
				// no-op
			}
		}()
	}
	wg.Wait()
	return nil
}

func getUniqueFilename(dir string, filename string) (string, error) {
	var fileExt string

	matches := duplicateFilenamePattern.FindStringSubmatch(filename)
	if len(matches) == 4 {
		filename = matches[1]
		fileExt = matches[3]
	} else {
		fileExt = filepath.Ext(filename)
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	largest := getLargestInteger(dirEntries, filename)
	return formatPackDataFilename(filename, largest, fileExt), nil
}

func formatPackDataFilename(filename string, largest int, fileExt string) string {
	return fmt.Sprintf("%s-%d%s", filename, largest, fileExt)
}

func getLargestInteger(readDir []os.DirEntry, filename string) int {
	largest := 0
	for _, entry := range readDir {
		if entry.IsDir() {
			continue
		} else if strings.HasPrefix(entry.Name(), filename) {
			submatch := duplicateFilenamePattern.FindStringSubmatch(entry.Name())
			if submatch == nil {
				continue
			} else if value, err := strconv.Atoi(submatch[2]); err == nil && value > largest {
				largest = value
			}
		}
	}
	return largest
}
