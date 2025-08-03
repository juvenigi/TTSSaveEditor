package main

import (
	"embed"
	"tts-cache-manager-cli/app"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:generate go run ./util/script/fetch-img.go
func main() {
	if err := app.RunWails(assets); err != nil {
		panic(err)
	}
}
