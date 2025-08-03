package util

import (
	"testing"
	"time"
)

func Test_tmp(t *testing.T) {
	errCh := make(chan error)
	defer close(errCh)
	DeadGopherChannel(errCh)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second * 1):
		return
	}
}

func Test_Panic(t *testing.T) {
	errCh := make(chan error)
	defer close(errCh)

	go func(errCh chan error) {
		defer DeadGopherChannel(errCh)
		panic("oh no")
	}(errCh)

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal(err)
		} else {
			t.Log(err)
		}
	case <-time.After(time.Second * 1):
		return
	}
}
