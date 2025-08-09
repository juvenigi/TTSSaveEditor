package internal

import (
	"os"
	"testing"
)

func TestCacheManagerApi_WriteToPackData(t *testing.T) {
	t.Log(os.Executable())
	api := NewCacheManagerApi()
	api.Startup(t.Context())

	gameDir := api.propertiesController.GetApplicationProperties().GameDir
	if err := api.WriteToPackData(gameDir + "/TS_Save_2.json"); err != nil {
		t.Fatal(err)
	}
}
