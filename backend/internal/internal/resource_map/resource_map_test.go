package resource_map

import (
	"io/fs"
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

type MockEntry struct {
	isDir bool
	name  string
}

func (m *MockEntry) Name() string {
	return m.name
}

func (m *MockEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

func (m *MockEntry) Type() fs.FileMode {
	return 0
}
func (m *MockEntry) IsDir() bool {
	return m.isDir
}
