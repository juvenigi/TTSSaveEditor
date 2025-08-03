package tabletop_save

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

const fileThatExists = "../../../resources/test/pillow-test.json"
const cacheFileThatExists = "../../../resources/test/cache/pillowtestjson.json"

func TestLoadCacheStatus_FileExists(t *testing.T) {
	fileAbsFilepath, err := filepath.Abs(fileThatExists)
	cacheAbsFilepath, err := filepath.Abs(cacheFileThatExists)
	cacheAbsFilepath = strings.TrimSuffix(cacheAbsFilepath, filepath.Ext(cacheAbsFilepath))
	if err != nil {
		t.Fatal(err)
	}
	cacheDir := filepath.Dir(cacheFileThatExists)

	cache, err := NewCacheScanner(cacheDir)
	if err != nil {
		t.Fatal(err)
	}

	res := GameResource{
		jsonPath:    "not.important.here",
		resourceUrl: fileAbsFilepath,
		status:      Remote,
	}

	err = cache.doCheckIfCached(&res)
	if err != nil {
		t.Fatal(err)
	}

	fileScanned, fileOk := cache.resourceCached[filepath.Base(cacheAbsFilepath)]
	if !fileScanned || !fileOk {
		t.Fatal("File not found in cache")
	}
	if res.status != RemoteCached {
		t.Fatal("resource status not set to RemoteCached")
	}
}

func TestLoadCacheStatus_CacheEntryExists(t *testing.T) {
	const fictionalFile = "/home/user/joghourt/tasty.jpg"

	// assuming the cached file is in the same dir as main file
	file := filepath.Base(fictionalFile)
	cacheFilename := GetCacheFilename(file)

	cache := GameCache{resourceCached: make(map[string]bool)}
	cache.resourceCached[cacheFilename] = true

	res := GameResource{
		jsonPath:    "not.important.here",
		resourceUrl: fictionalFile,
		status:      0,
	}

	err := cache.doCheckIfCached(&res)
	if err != nil {
		t.Fatal(err)
	}
	fileScanned, fileOk := cache.resourceCached[cacheFilename]
	if !fileScanned || !fileOk {
		t.Fatal("File not found in cache")
	}
	if res.status != RemoteCached {
		t.Fatal("resource status not set to RemoteCached")
	}
}

func TestGetCachedResourceLoc(t *testing.T) {
	gameDir := filepath.Dir(fileThatExists)
	loc := getCachedResourceLoc(gameDir, fileThatExists)

	rel, err := filepath.Rel(loc, cacheFileThatExists)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rel)
	if strings.Compare(rel, ".") != 0 {
		t.Fatal("resource location not found")
	}
}

func TestSearchingSaveCache(t *testing.T) {
	fn := visitGen(t)
	dir := filepath.Dir(filepath.Dir(fileThatExists))

	if err := filepath.WalkDir(dir, fn); err != nil {
		t.Fatal(err)
	}
}

func visitGen(t *testing.T) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		return visit(path, d, err, t)
	}
}

func visit(path string, d fs.DirEntry, err error, t *testing.T) error {
	if err != nil {
		return err
	}
	t.Log(" ", path, d.IsDir())
	return nil
}
