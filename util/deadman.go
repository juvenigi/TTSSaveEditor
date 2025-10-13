package util

import (
	"fmt"
	"runtime/debug"
)

// DeadGopherChannel dead gophers tell no tales, unless you give them an error channel
func DeadGopherChannel(errCh chan<- error) {
	if r := recover(); r != nil {
		errCh <- fmt.Errorf("panic in goroutine: %v\n%s", r, debug.Stack())
	}
}
