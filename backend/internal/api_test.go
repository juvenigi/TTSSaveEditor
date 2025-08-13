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
	if err := api.WriteToPackData(filepath.Join(gameDir, "Saves", "TS_Save_2.json")); err != nil {
		t.Fatal(err)
	}
}
