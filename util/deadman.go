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

// input: context and a Supplier<Tuple<any,error>>
// problem: primitives and non-nil value returns if err != nil
func HotRunAsync[T any](fn func() (T, error), resCh chan<- T, errCh chan<- error) {
	go func() {
		defer DeadGopherChannel(errCh)

		if res, err := fn(); err != nil {
			errCh <- err
		} else {
			resCh <- res
		}
	}()
}
