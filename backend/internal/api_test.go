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
	if err := api.WriteToPackData(filepath.Join(gameDir, "Saves", "TS_Save_438.json"), "Customize me"); err != nil {
		t.Fatal(err)
	}
}

func TestCacheManagerApi_ConstructPackedSave(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	_ = api.propertiesController.GetApplicationProperties().GameDir
	packDir := api.propertiesController.GetApplicationProperties().PackDataDir
	//if err := api.WriteToPackData(filepath.Join(gameDir, "Saves", "TS_Save_2.json")); err != nil {
	//	t.Fatal(err)
	//}
	if err := api.ConstructPackedSave(filepath.Join(packDir, "TS_Save_438.pack.json"), "Custom Save Name Test"); err != nil {
		t.Fatal(err)
	}
}
