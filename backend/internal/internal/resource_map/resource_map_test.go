package resource_map

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPatternMatching(t *testing.T) {
	if duplicateFilenamePattern.FindStringSubmatch("9") != nil {
		t.Errorf("duplicateFilenamePattern.FindStringSubmatch returned %s", duplicateFilenamePattern.FindStringSubmatch("9"))
	}

	if submatch := duplicateFilenamePattern.FindStringSubmatch("test-9.xml"); submatch == nil {
		t.Fatalf("duplicateFilenamePattern.FindStringSubmatch returned %s", duplicateFilenamePattern.FindStringSubmatch("9"))
	} else if len(submatch) != 4 {
		t.Fatalf("duplicateFilenamePattern.FindStringSubmatch returned %d submatches", len(submatch))
	}

	if submatch := duplicateFilenamePattern.FindStringSubmatch("test-0"); submatch == nil {
		t.Fatalf("duplicateFilenamePattern.FindStringSubmatch returned %s", duplicateFilenamePattern.FindStringSubmatch("9"))
	} else if len(submatch) != 4 {
		t.Fatalf("duplicateFilenamePattern.FindStringSubmatch returned %d submatches", len(submatch))
	}

}

func TestCreateBlankPackDataJson(t *testing.T) {
	loc, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	jsonLoc := filepath.Join(filepath.Dir(loc), PackDataJson)

	err = CreateBlankPackDataJson(jsonLoc)
	if err != nil {
		t.Fatal(err)
	}

	pd, err := InitPackDataFromFile(jsonLoc)
	if err != nil {
		t.Fatal(err)
	}
	if pd.packDataDir == "" {
		t.Fatal("pd is nil")
	}

}
