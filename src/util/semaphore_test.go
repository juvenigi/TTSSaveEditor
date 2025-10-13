package util

import (
	"sync"
	"testing"
	"time"
)

func TestSemaphore_Release(t *testing.T) {
	sem := NewSemaphore(1)
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	wg.Add(1)
	wg2.Add(2)

	go func() {
		sem.Acquire()
		wg.Wait()
		time.Sleep(1 * time.Second)
		sem.Release()
		wg2.Done()
	}()

	go func() {
		sem.Acquire()
		t.Log("second goroutine")
		wg2.Done()
	}()
	wg.Done()
	wg2.Wait()
	t.Log("test successful")
}
