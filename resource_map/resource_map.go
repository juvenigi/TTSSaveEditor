package resource_map

import (
	"crypto/sha3"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"tts-cache-manager-cli/httpfetcher"
)

type yamlPackData struct {
	Sha3Archive map[string]string `yaml:"sha-3-map"`
	UrlArchive  map[string]string `yaml:"url-map"`
}

type PackData struct {
	packDataDir   string
	packDataFiles map[string]bool
	saveTemplates []string
	sha3Archive   map[string]string
	urlArchive    map[string]string
}

func NewPackData(resourceMapFile string) (PackData, error) {
	var result PackData
	var yamlPd yamlPackData
	result.packDataDir = filepath.Dir(resourceMapFile)

	src, err := os.Open(resourceMapFile)
	if err != nil {
		return result, err
	}

	if entries, err := os.ReadDir(result.packDataDir); err != nil {
		return result, err
	} else {
		result.saveTemplates = make([]string, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasPrefix(entry.Name(), ".pack.json") {
				result.saveTemplates = append(result.saveTemplates, entry.Name())
			} else {
				result.packDataFiles[entry.Name()] = true
			}
		}
	}

	if err = yaml.NewDecoder(src).Decode(&yamlPd); err != nil {
		return result, err
	}
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
	return result, err
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

func (rm *PackData) AddLocalFile(url string, hasher *sha3.SHA3) (string, error) {
	var destFilename string
	blob, err := os.ReadFile(url)
	if err != nil {
		return destFilename, err
	}
	hasher.Reset()
	sum := hasher.Sum(blob)
	if alreadyPresent, ok := rm.sha3Archive[string(sum)]; ok {
		return alreadyPresent, nil
	}

	filename := filepath.Base(url)
	if rm.packDataFiles[filename] == true {
		destFilename, err = getUniqueFilename(rm.packDataDir, filename)
		if err != nil {
			return destFilename, err
		}
	} else {
		destFilename = filename
	}

	if err := os.WriteFile(destFilename, blob, os.ModePerm); err != nil {
		return destFilename, err
	}
	return destFilename, nil
}

func (rm *PackData) AddRemoteFile(url string, hasher *sha3.SHA3) error {
	res, err := httpfetcher.GetResource(url)
	if err != nil {
		return err
	}
	deduplicatedFilename, err := getUniqueFilename(rm.packDataDir, filepath.Base(url))
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(rm.packDataDir, deduplicatedFilename), res.Data, os.ModePerm); err != nil {
		return err
	}
	rm.packDataFiles[deduplicatedFilename] = true
	hasher.Reset()
	rm.sha3Archive[string(hasher.Sum(res.Data))] = deduplicatedFilename
	rm.urlArchive[url] = deduplicatedFilename

	return nil
}

func (rm *PackData) AddRemoteCachedFile(url string, cachedLoc string, hasher *sha3.SHA3) error {
	destFilename, err := rm.AddLocalFile(cachedLoc, hasher)
	if err != nil {
		return err
	}

	rm.urlArchive[url] = destFilename
	return nil
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
	deletedOnce := false
	deletedOnce = cleanupMap(rm.scanMap(rm.urlArchive), rm.urlArchive)
	deletedOnce = cleanupMap(rm.scanMap(rm.sha3Archive), rm.sha3Archive)

	var err error
	if deletedOnce {
		err = rm.FlushToDisk()
	}
	return err
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

func (rm *PackData) scanMap(urlToFileMap map[string]string) map[string]presentKey {
	var mappedFileExists = make(map[string]presentKey)
	for url, file := range urlToFileMap {
		if _, ok := rm.packDataFiles[file]; ok {
			mappedFileExists[file] = presentKey{present: true, key: url}
		} else {
			mappedFileExists[file] = presentKey{present: false, key: url}
		}
	}
	return mappedFileExists
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
	return fmt.Sprintf("%s-%d%s", filename, largest, fileExt), nil
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
