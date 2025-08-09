package internal

import (
	"strings"
	"testing"
)

func TestDropBoxFile(t *testing.T) {
	given := "https://www.dropbox.com/s/36xemw0ahjgq1lp/InvincibleReason.png?dl=1"
	expected := "httpswwwdropboxcoms36xemw0ahjgq1lpInvincibleReasonpngdl1"
	actual := GetCacheFilename(given)

	if strings.Compare(actual, expected) != 0 {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestGetEncodedString(t *testing.T) {
	given := "https://example.com/file name.png"

	actual := GetCacheFilename(given)

	expected := "httpsexamplecomfilenamepng"
	if strings.Compare(actual, expected) != 0 {
		t.Fatal("Expected", expected, "Got", actual)
	}
}

func TestGetEncodedString_NoExtensionFile(t *testing.T) {
	given := "https://steamusercontent-a.akamaihd.net/ugc/1752434872401823536/686EC0759C08456A15476928722EDDF8DE64AB0F/"

	actual := GetCacheFilename(given)

	expected := "httpssteamusercontentaakamaihdnetugc1752434872401823536686EC0759C08456A15476928722EDDF8DE64AB0F"
	if strings.Compare(actual, expected) != 0 {
		t.Fatal("Expected", expected, "Got", actual)
	}
}
