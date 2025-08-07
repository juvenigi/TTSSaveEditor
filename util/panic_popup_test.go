package util

import "testing"

// note: you must click ok for panic to proceed, if I didn't block it, then the window would've auto-closed
func TestShowWindowOnPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	defer ShowWindowOnPanic()
	panic("oh no")
}
