package httpfetcher

import (
	"bytes"
	"crypto/sha3"
	"io"
	"os"
	"testing"
)

func TestUrlResource_Checksum(t *testing.T) {
	given, err := os.Open("../../../resources/img/mistborn.jpg")
	if err != nil {
		t.Fatalf("failed to open file: %s", err.Error())
	}
	defer given.Close()
	allBytes, err := io.ReadAll(given)
	if err != nil {
		t.Fatalf("failed to read file: %s", err.Error())
	}

	new256 := sha3.New256()
	_, _ = new256.Write(allBytes)
	expected := new256.Sum(nil)

	actualRes, err := GetResource("https://cdnb.artstation.com/p/assets/panos/panos/006/685/241/large/a2a289d664d2cde2.jpg")
	if err != nil {
		t.Fatalf("failed to calculate checksum: %s", err.Error())
	}
	if actualRes.Data == nil {
		t.Fatalf("actual is nil")
	}
	if !bytes.Equal(expected, actualRes.checksum) {
		t.Fatalf("UrlResource.Checksum returned a wrong checksum")
	}
}
