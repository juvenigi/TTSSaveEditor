package resource_map

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/wrapped_io"
)

const PackDataJson = "resource-map.json"

type packDataJson struct {
	UrlArchive map[string]string `json:"url-map"`
}

type PackData struct {
	PackDataDir     string
	PackDataYamlLoc string
	Sha3Archive     map[string]string
	UrlArchive      map[string]string
}

func CreateBlankPackDataJson(packDataLocation string) error {
	var jsonD packDataJson
	file, err := os.OpenFile(packDataLocation, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil
		}
		return err
	}
	defer file.Close()
	if err = json.NewEncoder(file).Encode(jsonD); err != nil {
		return err
	}

	return nil
}

func InitPackDataFromFile(resourceMapFile string) (PackData, error) {
	var alreadyClosed bool
	var result PackData
	var jsonPd packDataJson
	result.PackDataDir = filepath.Dir(resourceMapFile)
	result.PackDataYamlLoc = resourceMapFile
	result.Sha3Archive = make(map[string]string)
	result.UrlArchive = make(map[string]string)

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
	if err = json.NewDecoder(src).Decode(&jsonPd); err != nil {
		return result, err
	}
	_ = src.Close()
	alreadyClosed = true

	if jsonPd.UrlArchive != nil {
		result.UrlArchive = jsonPd.UrlArchive
	}

	entries, err := os.ReadDir(result.PackDataDir)
	if err != nil {
		return result, err
	}

	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		_, sum := wrapped_io.GetFileAndChecksum(filepath.Join(result.PackDataDir, entry.Name()))
		result.Sha3Archive[hex.EncodeToString(sum)] = entry.Name()
	}

	if err = result.cleanup(); err != nil {
		return result, err
	}

	return result, nil
}

func (rm *PackData) FlushToDisk() error {
	file, err := os.Create(filepath.Join(rm.PackDataYamlLoc))
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(packDataJson{
		UrlArchive: rm.UrlArchive,
	}); err != nil {
		return err
	}

	return file.Sync()
}

func (rm *PackData) cleanup() error {
	filenames, err := NewDirectorySnapshot(rm.PackDataDir)
	if err != nil {
		return err
	}

	deletedOnce := false
	urlScanMap := scanMap(rm.UrlArchive, filenames)
	deletedOnce = cleanupMap(urlScanMap, &rm.UrlArchive)

	if deletedOnce {
		return rm.FlushToDisk()
	}

	return nil
}

// DirectorySnapshot exists to avoid repeatedly calling os.ReadDir every time we add a new file to PackData.
type DirectorySnapshot struct {
	Entries map[string]struct{}
}

func NewDirectorySnapshot(packDataDir string) (*DirectorySnapshot, error) {
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

	return &DirectorySnapshot{
		Entries: filenames,
	}, nil
}

func (filenames *DirectorySnapshot) GetUniqueFilename(unwritten string) string {
	attempt := 0
	attemptedName := fmt.Sprintf("%s-%d", unwritten, attempt)
	// try without attempt counter first, then iterate over attempts
	if _, present := filenames.Entries[unwritten]; !present {
		return unwritten
	}
	for {
		// loop until a unique filename is found
		if _, present := filenames.Entries[attemptedName]; !present {
			return attemptedName
		}
		attempt++
	}
}

func scanMap(urlToFileMap map[string]string, filenames *DirectorySnapshot) map[string]presentKey {
	var mappedFileExists = make(map[string]presentKey)
	for url, file := range urlToFileMap {
		if _, ok := filenames.Entries[file]; ok {
			mappedFileExists[file] = presentKey{present: true, key: url}
		} else {
			mappedFileExists[file] = presentKey{present: false, key: url}
		}
	}
	return mappedFileExists
}

// filename, count, extension(may be empty)
var duplicateFilenamePattern = regexp.MustCompile("(.*)-(\\d+)(\\.?.*)$")

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

// WriteSaveToPackData note: this is not thread safe
func (rm *PackData) WriteSaveToPackData(originalName string, blob []byte) error {
	trimmedName := strings.TrimSuffix(originalName, ".json")
	candidateName := trimmedName

	dirEntries, err := os.ReadDir(rm.PackDataDir)
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
		if err = os.WriteFile(filepath.Join(rm.PackDataDir, candidateName+".pack.json"), blob, 0644); err != nil {
			return err
		}
		return nil
	}

	counter := 0
	for {
		candidateName = fmt.Sprintf("%s-%d", trimmedName, counter)
		if _, ok := nameSet[candidateName+".pack.json"]; !ok {
			if err = os.WriteFile(filepath.Join(rm.PackDataDir, candidateName+".pack.json"), blob, 0644); err != nil {
				return err
			}
			return nil
		}
		counter++
	}
}

func (rm *PackData) GetAllFilesInPackData() ([]string, error) {
	var packJsonLocs []string
	entries, err := os.ReadDir(rm.PackDataDir)
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
			dest := filepath.Join(rm.PackDataDir, entry.Name())
			if err := wrapped_io.CopyFile(src, dest); err != nil {
				return err
			}
		}
	}

	foreignPd, err := InitPackDataFromFile(filepath.Join(otherDir, PackDataJson))
	if err != nil {
		return err
	}

	maps.Copy(rm.UrlArchive, foreignPd.UrlArchive)
	maps.Copy(rm.Sha3Archive, foreignPd.Sha3Archive)

	return rm.FlushToDisk()
}

var bakNumPattern = regexp.MustCompile(`\d+`)

func (rm *PackData) MakeBackup() error {
	parent, dir := filepath.Split(rm.PackDataDir)
	backupDir, err := getBackupDirName(parent, dir)
	if err != nil {
		return err
	}
	if err = os.Mkdir(backupDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(rm.PackDataDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(rm.PackDataDir, entry.Name())
		dest := filepath.Join(backupDir, entry.Name())
		if err := wrapped_io.CopyFile(src, dest); err != nil {
			return err
		}
	}

	return nil
}

func getBackupDirName(parent string, dir string) (string, error) {
	dirEntries, err := os.ReadDir(parent)
	if err != nil {
		return "", err
	}
	var packs []string
	largest := -1
	for _, entry := range dirEntries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), dir) {
			packs = append(packs, entry.Name())
			if digits := bakNumPattern.FindString(entry.Name()); len(digits) > 0 {
				cursor, err := strconv.Atoi(digits)
				if err != nil {
					return "", err
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

	return filepath.Join(parent, fmt.Sprintf("%s.%d", strippedDir, largest+1)), nil
}
