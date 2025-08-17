package util

import (
	"fmt"
	"runtime"

	"github.com/ncruces/zenity"
)

func ShowWindowOnPanic() {
	if r := recover(); r != nil {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		stackTrace := string(buf[:n])

		msg := fmt.Sprintf("Application Error:\n%v\n\n%s", r, stackTrace)

		// must block here, otherwise the user won't even see the window or see a window open and close
		_ = zenity.Error(msg, zenity.Title("Unexpected Error"), zenity.ErrorIcon)

		panic(r)
	}
}
