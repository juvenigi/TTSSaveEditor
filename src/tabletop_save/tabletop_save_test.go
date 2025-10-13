package tabletop_save

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTabletopSave_Init(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	fixturePath := filepath.Join(wd, "..", "..", "..", "test", "fixtures", "TS_Save_71.json")

	var save TabletopSave
	err = save.Init(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(save.resources) != 4 {
		t.Fatalf("wrong number of resources, expected 4, got %d", len(save.resources))
	}
}
