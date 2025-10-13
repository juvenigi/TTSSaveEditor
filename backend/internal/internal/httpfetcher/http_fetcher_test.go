package httpfetcher

import (
	"bytes"
	"crypto/sha3"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestUrlResource_Checksum(t *testing.T) {
	given, err := os.Open("../../../../resources/img/mistborn.jpg")
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

	actualRes, err := GetResourceAsync(t.Context(), http.DefaultClient, "https://cdnb.artstation.com/p/assets/panos/panos/006/685/241/large/a2a289d664d2cde2.jpg")
	if err != nil {
		t.Fatalf("failed to calculate Checksum: %s", err.Error())
	}
	if actualRes.Data == nil {
		t.Fatalf("actual is nil")
	}
	if !bytes.Equal(expected, actualRes.Checksum) {
		t.Fatalf("UrlResource.Checksum returned a wrong Checksum")
	}
}

var testErr = errors.New("test error")

func TestError_1(t *testing.T) {
	another := errors.New("another error")
	combined := fmt.Errorf("%w: %w", testErr, another)
	t.Log(combined)
	e := combined.(interface{ Unwrap() []error })
	t.Log(e.Unwrap())
	if !errors.Is(combined, testErr) {
		t.Fatalf("combined error does not match")
	}

}

func TestError_2(t *testing.T) {
	another := errors.New("another error")
	combined := errors.Join(testErr, another)
	t.Log(combined)
	e := combined.(interface{ Unwrap() []error })
	t.Log(e.Unwrap())
	if !errors.Is(combined, testErr) {
		t.Fatalf("combined error does not match")
	}

}

func TestTime(t *testing.T) {
	zero := time.Time{}
	t.Log(zero.GoString())
	if zero.After(time.Now()) {
		t.Fatalf("zero time is in the future")
	}
}
