package util

import (
	"strings"
	"testing"
)

func TestLastSegment(t *testing.T) {
	given := "foo.bar.baz"
	expected := "baz"
	actual := LastSegment(given)
	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("LastSegment(%q): expected %q, actual %q", given, expected, actual)
	}
}
