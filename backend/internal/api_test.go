package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheManagerApi_WriteToPackData(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	gameDir := api.propertiesController.GetApplicationProperties().GameDir
	if err := api.WriteToPackData(filepath.Join(gameDir, "Saves", "TS_Save_68.json"), "Hmm we have a bug"); err != nil {
		t.Fatal(err)
	}
}

func TestCacheManagerApi_ConstructPackedSave(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	_ = api.propertiesController.GetApplicationProperties().GameDir
	packDir := api.propertiesController.GetApplicationProperties().PackDataDir
	if err := api.ConstructPackedSave(filepath.Join(packDir, "TS_Save_1.pack.json"), "Serrogo testo"); err != nil {
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
