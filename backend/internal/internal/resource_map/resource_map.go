package resource_map

import (
	"bytes"
	"context"
	"crypto/sha3"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"tts-cache-manager-cli/backend/internal/internal/httpfetcher"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save"
	"tts-cache-manager-cli/backend/internal/internal/wrapped_io"

	"gopkg.in/yaml.v3"
)

const PackDataYaml = "resource-map.yaml"

type yamlPackData struct {
	UrlArchive map[string]string `yaml:"url-map"`
}

type PackData struct {
	packDataDir     string
	packDataYamlLoc string
	sha3Archive     map[string]string
	urlArchive      map[string]string
}

func CreateBlankPackDataYaml(packDataLocation string) error {
	var yamlData yamlPackData
	file, err := os.OpenFile(packDataLocation, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil
		}
		return err
	}
	defer file.Close()
	if err = yaml.NewEncoder(file).Encode(yamlData); err != nil {
		return err
	}

	return nil
}

func InitPackDataFromFile(resourceMapFile string) (PackData, error) {
	var alreadyClosed bool
	var result PackData
	var yamlPd yamlPackData
	result.packDataDir = filepath.Dir(resourceMapFile)
	result.packDataYamlLoc = resourceMapFile
	result.sha3Archive = make(map[string]string)
	result.urlArchive = make(map[string]string)

	src, err := os.Open(resourceMapFile)
	if err != nil {
		return result, err
	}
	defer func(src *os.File, closed bool) {
		if !closed {
			_ = src.Close()
		}
	}(src, alreadyClosed)

	// todo: could one use a better solution here? (invalid yaml gets skipped without the user knowing)
	if err = yaml.NewDecoder(src).Decode(&yamlPd); err != nil {
		return result, err
	}
	_ = src.Close()
	alreadyClosed = true

	if yamlPd.UrlArchive != nil {
		result.urlArchive = yamlPd.UrlArchive
	}

	entries, err := os.ReadDir(result.packDataDir)
	if err != nil {
		return result, err
	}

	extensions := []string{".yaml", ".yml", ".json"}
	for _, entry := range entries {
		if entry.IsDir() || slices.Contains(extensions, filepath.Ext(entry.Name())) {
			continue
		}
		_, sum := wrapped_io.GetFileAndChecksum(filepath.Join(result.packDataDir, entry.Name()))
		result.sha3Archive[entry.Name()] = hex.EncodeToString(sum)
	}

	if err = result.cleanup(); err != nil {
		return result, err
	}

	return result, nil
}

func (rm *PackData) FlushToDisk() error {
	file, err := os.Create(filepath.Join(rm.packDataYamlLoc))
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewEncoder(file).Encode(yamlPackData{
		UrlArchive: rm.urlArchive,
	}); err != nil {
		return err
	}

	return file.Sync()
}

// filename, count, extension(may be empty)
var duplicateFilenamePattern = regexp.MustCompile("(.*)-(\\d+)(\\.?.*)$")

// RwFilenameSlice exists to avoid repeatedly calling os.ReadDir every time we add a new file to PackData.
type RwFilenameSlice struct {
	mu      *sync.RWMutex
	entries map[string]struct{}
}

func NewRwFilenameSlice(packDataDir string) (*RwFilenameSlice, error) {
	fileEntries, err := os.ReadDir(packDataDir)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func (rm *PackData) AddLocalFile(res *tabletop_save.GameResource, filenames *RwFilenameSlice) error {
	var location = strings.TrimPrefix(res.ResourceUrl, "file:///")
	blob, sum := wrapped_io.GetFileAndChecksum(location)
	if sum == nil {
		return errors.New("file not found")
	}

	packDataFile, exists := rm.sha3Archive[hex.EncodeToString(sum)]
	if exists {
		res.ResourceUrl = "pack:///" + packDataFile
		res.Status = tabletop_save.Packed
		return nil
	}

	var destFilename = tabletop_save.GetCacheFilename(filepath.Base(res.ResourceUrl))
	return rm.writeResource(filenames, destFilename, blob, res)
}

func (rm *PackData) AddRemoteFile(ctx context.Context, res *tabletop_save.GameResource, filenames *RwFilenameSlice) error {
	urlRes, err := httpfetcher.GetResourceAsync(ctx, res.ResourceUrl)
	if err != nil {
		res.Status = tabletop_save.Failed
		return err
	}

	packFileName, ok := rm.sha3Archive[hex.EncodeToString(urlRes.Checksum)]
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
	return rm.writeResource(filenames, escapedName, urlRes.Data, res)
}

func (rm *PackData) writeResource(filenames *RwFilenameSlice, destinationName string, blob []byte, res *tabletop_save.GameResource) error {
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

	return rm.flushAndUpdate(res, blob, filenames)
}

func (rm *PackData) flushAndUpdate(res *tabletop_save.GameResource, file []byte, filenames *RwFilenameSlice) error {
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

	new256 := sha3.New256()
	if _, criticalErr = new256.Write(file); criticalErr != nil {
		return criticalErr
	}

	rm.urlArchive[res.ResourceUrl] = written
	rm.sha3Archive[hex.EncodeToString(new256.Sum(nil))] = written

	res.ResourceUrl = "pack://" + written
	res.Status = tabletop_save.Packed

	return nil
}

func attemptToWriteExtLess(dir string, name string, count int, content []byte) (string, []string, error) {
	tries := 0
	var previousAttempts []string
	attemptName := name
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

func (rm *PackData) AddRemoteCachedFile(res *tabletop_save.GameResource, filenames *RwFilenameSlice, scanner *tabletop_save.GameCacheFinder) error {
	fpath, ok := scanner.ResourceCached[tabletop_save.GetCacheFilename(res.ResourceUrl)]
	if !ok {
		return errors.New("cached file not found")
	}

	file, fileChecksum := wrapped_io.GetFileAndChecksum(fpath)
	if fileChecksum == nil {
		return errors.New("cached file not found")
	}

	packDatafile, ok := rm.sha3Archive[hex.EncodeToString(fileChecksum)]
	if ok {
		res.ResourceUrl = "pack://" + packDatafile
		res.Status = tabletop_save.Packed
		return nil
	} else {
		destinationFilename := tabletop_save.GetCacheFilename(res.ResourceUrl)
		return rm.writeResource(filenames, destinationFilename, file, res)
	}
}

func (rm *PackData) cleanup() error {
	filenames, err := NewRwFilenameSlice(rm.packDataDir)
	if err != nil {
		return err
	}

	deletedOnce := false
	urlScanMap := rm.scanMap(rm.urlArchive, filenames)
	deletedOnce = cleanupMap(urlScanMap, &rm.urlArchive)

	if deletedOnce {
		return rm.FlushToDisk()
	}

	return nil
}

type presentKey struct {
	present bool
	key     string
}

// todo: proper test for cleanup, validation, init
func cleanupMap(mappedUrlExists map[string]presentKey, target *map[string]string) bool {
	deletedOnce := false
	if target == nil {
		return false
	}

	for _, key := range mappedUrlExists {
		if !key.present {
			delete(*target, key.key)
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

func (rm *PackData) ImportFromSave(ctx context.Context, data *tabletop_save.TSSaveFile, gameDir string) error {
	resources := data.GetAllResources()
	cacheFinder, err := tabletop_save.InitGameCacheFinder(gameDir)
	if err != nil {
		return err
	}

	filenames, err := NewRwFilenameSlice(rm.packDataDir)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		cacheFinder.MutCacheStatus(resource)
	}
	// may contain filenames or urls
	var seen = make(map[string]struct{})
	for idx := range resources {
		var resourceRef = resources[idx]
		if _, already := seen[resourceRef.ResourceUrl]; already {
			continue
		}
		seen[resourceRef.ResourceUrl] = struct{}{}

		var errSw error
		switch resourceRef.Status {
		case tabletop_save.Remote, tabletop_save.RemoteCached:
			if packUrl, ok := rm.urlArchive[resourceRef.ResourceUrl]; ok {
				resourceRef.ResourceUrl = "pack://" + packUrl
				resourceRef.Status = tabletop_save.Packed
			} else if resourceRef.Status == tabletop_save.Remote {
				errSw = rm.AddRemoteFile(ctx, resourceRef, filenames)
			} else {
				errSw = rm.AddRemoteCachedFile(resourceRef, filenames, &cacheFinder)
			}
		case tabletop_save.Local:
			errSw = rm.AddLocalFile(resourceRef, filenames)
		default:
			// no-op
		}
		if errSw != nil {
			log.Println(errSw)
			errSw = nil
		}
	}

	return rm.FlushToDisk()
}

// note: this is not thread safe
func (rm *PackData) WriteSaveToPackData(originalName string, blob []byte) error {
	trimmedName := strings.TrimSuffix(originalName, ".json")
	candidateName := trimmedName

	dirEntries, err := os.ReadDir(rm.packDataDir)
	if err != nil {
		return err
	}
	nameSet := make(map[string]struct{})
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pack.json") {
			nameSet[entry.Name()] = struct{}{}
		}
	}
	if _, ok := nameSet[candidateName+".pack.json"]; !ok {
		if err = os.WriteFile(filepath.Join(rm.packDataDir, candidateName+".pack.json"), blob, 0644); err != nil {
			return err
		}
		return nil
	}

	counter := 0
	for {
		candidateName = fmt.Sprintf("%s-%d", trimmedName, counter)
		if _, ok := nameSet[candidateName+".pack.json"]; !ok {
			if err = os.WriteFile(filepath.Join(rm.packDataDir, candidateName+".pack.json"), blob, 0644); err != nil {
				return err
			}
			return nil
		}
		counter++
	}
}

const packLen = len("pack://")

func (rm *PackData) LocalizePackedResources(resources []*tabletop_save.GameResource) error {
	packJsonLocs, err := getAllFilesInPackData(rm)
	if err != nil {
		return err
	}

	for idx := range resources {
		var res = resources[idx]
		if res.Status == tabletop_save.Packed && strings.HasPrefix(res.ResourceUrl, "pack://") {
			filename := res.ResourceUrl[packLen:]
			if slices.Contains(packJsonLocs, filename) {
				res.ResourceUrl = "file:///" + filepath.Join(rm.packDataDir, filename)
			} else {
				return fmt.Errorf("could not find resource %s in PackData", res.ResourceUrl)
			}
		} else if res.Status == tabletop_save.Remote {
			if packDataLoc, ok := rm.urlArchive[res.ResourceUrl]; ok {
				res.Status = tabletop_save.Packed
				res.ResourceUrl = "file:///" + filepath.Join(rm.packDataDir, packDataLoc)
			} else {
				continue
			}
		}
	}

	return nil
}

func getAllFilesInPackData(rm *PackData) ([]string, error) {
	var packJsonLocs []string
	entries, err := os.ReadDir(rm.packDataDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		// oh no, terrible inefficiency!
		if !entry.IsDir() && !strings.HasSuffix(entry.Name(), ".pack.json") && !strings.HasSuffix(entry.Name(), ".yaml") {
			packJsonLocs = append(packJsonLocs, entry.Name())
		}
	}
	return packJsonLocs, nil
}

func (rm *PackData) MergeWith(otherDir string, makeBackup bool) error {
	if makeBackup {
		if err := rm.MakeBackup(); err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(otherDir)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(entries, func(s os.DirEntry) bool {
		return !s.IsDir() && strings.HasSuffix(s.Name(), ".yaml")
	})
	if idx == -1 {
		return errors.New("could not find pack data index file")
	}
	entries = append(entries[:idx], entries[idx+1:]...)
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		} else if strings.HasSuffix(entry.Name(), ".pack.json") {
			src := filepath.Join(otherDir, entry.Name())
			if srcBlob, err := os.ReadFile(src); err != nil {
				return err
			} else {
				if err := rm.WriteSaveToPackData(entry.Name(), srcBlob); err != nil {
					return err
				}
			}
		} else {
			src := filepath.Join(otherDir, entry.Name())
			dest := filepath.Join(rm.packDataDir, entry.Name())
			if err := wrapped_io.CopyFile(src, dest); err != nil {
				return err
			}
		}
	}

	foreignPd, err := InitPackDataFromFile(filepath.Join(otherDir, PackDataYaml))
	if err != nil {
		return err
	}

	maps.Copy(rm.urlArchive, foreignPd.urlArchive)
	maps.Copy(rm.sha3Archive, foreignPd.sha3Archive)

	return rm.FlushToDisk()
}

var bakNumPattern = regexp.MustCompile(`\d+`)

func (rm *PackData) MakeBackup() error {
	parent, dir := filepath.Split(rm.packDataDir)
	dirEntries, err := os.ReadDir(parent)
	if err != nil {
		return err
	}
	var packs []string
	largest := -1
	for _, entry := range dirEntries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), dir) {
			packs = append(packs, entry.Name())
			if digits := bakNumPattern.FindString(entry.Name()); len(digits) > 0 {
				cursor, err := strconv.Atoi(digits)
				if err != nil {
					return err
				}
				if cursor > largest {
					largest = cursor
				}
			}
		}
	}
	var strippedDir string
	idx := strings.Index(dir, ".")
	if idx == -1 {
		strippedDir = dir
	} else {
		strippedDir = dir[:idx]
	}

	backupDir := filepath.Join(parent, fmt.Sprintf("%s.%d", strippedDir, largest+1))
	if err = os.Mkdir(backupDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(rm.packDataDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(rm.packDataDir, entry.Name())
		dest := filepath.Join(backupDir, entry.Name())
		if err := wrapped_io.CopyFile(src, dest); err != nil {
			return err
		}
	}

	return nil
}

func (rm *PackData) CachePackedResources(resList []*tabletop_save.GameResource, gameDir string) error {
	for _, res := range resList {
		if res.Status == tabletop_save.Packed {
			if err := copyPackedToGameCache(res, gameDir, rm.packDataDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyPackedToGameCache(res *tabletop_save.GameResource, gameDir string, packDataDir string) error {
	if res.Status != tabletop_save.Packed {
		return fmt.Errorf("unexpected resource type %d", res.Status)
	}

	subfolder := "unpacked"
	var lastSlash = strings.LastIndex(res.JsonPointer, "/")
	lastJsonPointerSegment := res.JsonPointer[lastSlash+1:]

	var filename = strings.TrimPrefix(res.ResourceUrl, "pack://")

	var gameCacheFolder string
	switch lastJsonPointerSegment {
	case "AssetBundleURL":
		gameCacheFolder = "Assetbundles"

	case "ImageURL", "ImageSecondaryURL", "DiffuseURL", "NormalURL", "SpecularURL", "SkyURL":
		gameCacheFolder = "Images"

	case "MeshURL", "ColliderURL":
		gameCacheFolder = "Models"

	case "AudioURL":
		gameCacheFolder = "Audio"

	case "XmlUIURL":
		gameCacheFolder = "UI"

	default:
		return fmt.Errorf("unknown json pointer: %s", lastJsonPointerSegment)
	}

	srcPath := filepath.Join(packDataDir, filename)

	finalPath := filepath.Join(gameDir, gameCacheFolder, subfolder, filename)

	if err := wrapped_io.CopyFile(srcPath, finalPath); err != nil {
		return err
	}

	res.ResourceUrl = fmt.Sprintf("%s/%s", subfolder, filename)

	return nil
}
