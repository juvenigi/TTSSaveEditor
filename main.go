package main

import (
	"TTSBundler/rest"
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	go func() {
		if err := rest.CreateSaveEditorBackend().Listen(":3000"); err != nil {
			panic(err)
		}
	}()
	app := NewApp(context.Background())
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wails-sandbox",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Bind:             []any{app},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
