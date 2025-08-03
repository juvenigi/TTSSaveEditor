package util

import (
	"errors"
	"strings"
	"testing"
)

func TestNewResolvers_Resolve(t *testing.T) {
	promise, resolve, _ := NewResolvers[string]()
	resolve("foo")

	actual, err := promise.Await(t.Context())
	if err != nil {
		t.Fatal(err)
	} else if strings.Compare(*actual, "foo") != 0 {
		t.Errorf("expecting foo, got %s", *actual)
	}
	t.Log(*actual)
}

func TestNewResolvers_Reject(t *testing.T) {
	promise, _, reject := NewResolvers[string]()
	reject(errors.New("error"))

	_, err := promise.Await(t.Context())
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
}
