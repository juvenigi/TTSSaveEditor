package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCacheManagerApi_WriteToPackData(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	gameDir := api.applicationProperties.GameDir
	if err := api.WriteToPackData(filepath.Join(gameDir, "Saves", "TS_Save_68.json"), "Hmm we have a bug"); err != nil {
		t.Fatal(err)
	}
}

func TestCacheManagerApi_ConstructPackedSave(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	packDir := api.applicationProperties.PackDataDir
	if err := api.CreateSingleplayerSave(filepath.Join(packDir, "TS_Save_68.pack.json"), "Serrogo testo"); err != nil {
		t.Fatal(err)
	}
}

func TestCacheManagerApi_WriteToPackDataWithCache(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	packDir := api.applicationProperties.PackDataDir
	if err := api.ConstructPackCachedSave(filepath.Join(packDir, "TS_Save_68.pack.json"), "willYouLoadBruh"); err != nil {
		t.Fatal(err)
	}
}

func TestCacheManagerApi_MergePackDataWithAnother(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("executable dir:", filepath.Dir(executable))
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	otherLoc := filepath.Join(filepath.Dir(executable), "NeoData")

	if err := api.MergePackDataWithAnother(otherLoc, true); err != nil {
		t.Fatal(err)
	}
}

// todo: write a proper test
func TestCacheManagerApi_OpenInFileExploererWithCache(t *testing.T) {
	exe := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator\\Saves\\"
	fromSlash := filepath.FromSlash(strings.ReplaceAll(filepath.Dir(exe), string(filepath.Separator), "/"))
	t.Log(fromSlash)
	t.Log(filepath.ToSlash(fromSlash))

	dir, _ := filepath.Split(exe)

	dirDir, file := filepath.Split(dir)

	t.Log("dir:", dir)
	t.Log("dirDir:", dirDir)
	t.Log("file:", len(file))

}
