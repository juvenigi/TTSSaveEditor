package main

import (
	"embed"
	"tts-cache-manager-cli/app"
	"tts-cache-manager-cli/util"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:generate go run ./util/script/fetch-img.go
func main() {
	defer util.ShowWindowOnPanic()
	if err := app.RunWails(assets); err != nil {
		panic(err)
	}
}
