package tabletop_save

import (
	"testing"
	"time"
	"tts-cache-manager-cli/properties"
)

func TestSaveManager_GetTabletopSaveFilesAsync(t *testing.T) {
	var filesCh = make(chan TSSaveFile)
	var errCh = make(chan error)
	defer close(filesCh)
	defer close(errCh)
	pp, err := properties.GetDefaultApplicationProperties()
	if err != nil {
		t.Fatal(err)
	}
	manager := NewSaveManager(pp)

	done, err := manager.GetTabletopSaveFilesAsync(filesCh, errCh)
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case save := <-filesCh:
			t.Log(save.savefileLocation)
		case err := <-errCh:
			t.Fatal(err)
		case <-done:
			t.Log("done")
			return
		case _ = <-time.After(1 * time.Second):
			t.Fatal("timeout")
		}
	}
}
