package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"testing"
)

func TestGetTabletopSaves(t *testing.T) {
	api := NewCacheManagerApi()

	go func() {
		runtime.EventsOn(t.Context(), NewFileEvent, func(data ...any) {
			t.Logf("len: %d\n", len(data))
			t.Log(data)
		})
	}()

	if err := api.GetTabletopSaves(); err != nil {
		t.Fatal(err)
	}
}
