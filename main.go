package main

import (
	"TTSBundler/grdive"
	ttssocket "TTSBundler/tts-socket"
	"context"
	"embed"
	"flag"
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.org/x/oauth2"
)

var (
	begForToken = flag.Bool("gdrive", false, "open auth consent window to obtain token?")
	editorAPI   = flag.Bool("editorapi", false, "listen to external editor API events?")
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	flag.Parse()
}

func main() {
	if *begForToken {
		if browser.OpenURL(grdive.AuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)) != nil {
			panic("cannot get API token!")
		}
	}
	if *editorAPI {
		go ttssocket.ListenToAppTCP(":39998")
	}
	//if err := rest.CreateSaveEditorBackend().Listen(":3000"); err != nil {
	//	log.Fatal(err)
	//}

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
