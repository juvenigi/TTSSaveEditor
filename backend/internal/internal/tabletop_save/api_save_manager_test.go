package tabletop_save

import (
	"testing"
	"time"
)

// todo: do not hardcode paths
func TestSaveManager_GetTabletopSaveFilesAsync(t *testing.T) {
	gameDir := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator"
	packDataDir := "C:\\Users\\joghourt\\IdeaProjects\\TabletopResourceManager\\build\\bin\\PackData"
	err, resChan := GetTabletopSaveFilesAsync(gameDir, packDataDir)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case save, ok := <-resChan:
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
		case <-time.After(10 * time.Second):
			t.Fatal("timeout")
		}
	}
}
