package properties

import (
	"path/filepath"
	"strings"
	"testing"
)

// the following tests for documentation purposes only

func TestPathSeparator(t *testing.T) {
	given := "C://foo/bar"
	actual := filepath.FromSlash(given)
	if strings.Compare(actual, "C:\\\\foo\\bar") != 0 {
		t.Fatal("Not what expected:", actual)
	}
}

func TestPathSeparatorIdentity(t *testing.T) {
	given := "C:\\\\foo\\bar"
	actual := filepath.FromSlash(given)
	if strings.Compare(actual, given) != 0 {
		t.Fatal("Expected no change, got:", actual)
	}

}
