package util

import (
	"fmt"
	"github.com/ncruces/zenity"
	"runtime"
)

func ShowWindowOnPanic() {
	if r := recover(); r != nil {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		stackTrace := string(buf[:n])

		msg := fmt.Sprintf("Application Error:\n%v\n\n%s", r, stackTrace)

		_ = zenity.Error(msg, zenity.Title("Unexpected Error"), zenity.ErrorIcon)

		panic(r)
	}
}
