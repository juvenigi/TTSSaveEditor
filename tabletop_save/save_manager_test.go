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

	err, filesCh := manager.GetTabletopSaveFilesAsync()
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case save, ok := <-filesCh:
			if !ok {

				t.Log("files channel closed")
				return
			} else {
				res := save.Result
				view := res.ToView()
				t.Log(view)
			}
			if save.Err != nil {
				t.Fatal(save.Err)
				return
			}

			t.Log(save.Result.savefileLocation)
		case <-time.After(199999999 * time.Second):
			t.Fatal("timeout")
		}
	}
}
