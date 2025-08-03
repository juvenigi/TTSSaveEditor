package tabletop_save

import (
	"testing"
	"time"
	"tts-cache-manager-cli/properties"
)

func TestSaveManager_GetTabletopSaveFilesAsync(t *testing.T) {
	pp, err := properties.GetDefaultApplicationProperties()
	if err != nil {
		t.Fatal(err)
	}
	manager := NewSaveManager(pp)

	err, filesCh, errCh := manager.GetTabletopSaveFilesAsync()
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case save, ok := <-filesCh:
			if !ok {
				t.Log("files channel closed")
				return
			}
			t.Log(save.savefileLocation)
		case err, ok := <-errCh:
			if !ok {
				t.Log("err channel closed")
				return
			}
			t.Fatal(err)
		case <-time.After(1 * time.Second):
			t.Fatal("timeout")
		}
	}
}
