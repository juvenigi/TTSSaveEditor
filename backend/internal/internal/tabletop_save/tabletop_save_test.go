package tabletop_save

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

const fileThatExists = "../../../resources/test/pillow-test.json"
const cacheFileThatExists = "../../../resources/test/cache/pillowtestjson.json"

func TestTSSaveFile_SaveAsPackTemplate(t *testing.T) {

	save := TSSaveFile{
		savefileLocation: "../resources/test/pillow-test.json",
		resourceBundle: []ResourceBundle{
			{
				JsonPointer: "",
				Name:        "",
				Nickname:    "",
				Guid:        "",
				resources: []GameResource{
					{
						JsonPointer: "/ObjectStates/0/CustomMesh/MeshURL",
						ResourceUrl: "packed://ohyes",
						Status:      Packed,
					},
				},
			},
		},
	}
	if _, err := save.GetPortableJsonBlob(nil); err != nil {
		t.Fatal(err)
	}
}

func TestLoadCacheStatus_FileExists(t *testing.T) {
	fileAbsFilepath, err := filepath.Abs(fileThatExists)
	cacheAbsFilepath, err := filepath.Abs(cacheFileThatExists)
	cacheAbsFilepath = strings.TrimSuffix(cacheAbsFilepath, filepath.Ext(cacheAbsFilepath))
	if err != nil {
		t.Fatal(err)
	}
	cacheDir := filepath.Dir(cacheFileThatExists)

	cache, err := InitGameCacheFinder(cacheDir)
	if err != nil {
		t.Fatal(err)
	}

	res := GameResource{
		JsonPointer: "not.important.here",
		ResourceUrl: fileAbsFilepath,
		Status:      Remote,
	}
	cache.MutCacheStatus(&res)
	//
	//fileScanned, fileOk := cache.ResourceCached[filepath.Base(cacheAbsFilepath)]
	//if !fileScanned || !fileOk {
	//	t.Fatal("File not found in cache")
	//}
	if res.Status != RemoteCached {
		t.Fatal("resource Status not set to RemoteCached")
	}
}

func TestLoadCacheStatus_CacheEntryExists(t *testing.T) {
	const fictionalFile = "/home/user/joghourt/tasty.jpg"

	// assuming the cached file is in the same dir as main file
	file := filepath.Base(fictionalFile)
	cacheFilename := GetCacheFilename(file)

	cache := GameCacheFinder{ResourceCached: make(map[string]string)}
	cache.ResourceCached[cacheFilename] = ""

	res := GameResource{
		JsonPointer: "not.important.here",
		ResourceUrl: fictionalFile,
		Status:      0,
	}

	//cache.MutCacheStatus(&res)
	//fileScanned, fileOk := cache.ResourceCached[cacheFilename]
	//if !fileScanned || !fileOk {
	//	t.Fatal("File not found in cache")
	//}
	if res.Status != RemoteCached {
		t.Fatal("resource Status not set to RemoteCached")
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

func TestRegexPattern(t *testing.T) {
	var foo = "$packrev://"
	t.Log(len(saveRevPattern.FindStringSubmatch(foo)))
}

func TestUnkownType(t *testing.T) {
	var foo any = 1

	var bar string

	bar = foo.(string)

	t.Log(bar)
}
