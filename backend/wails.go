package backend

import (
	"embed"
	"tts-cache-manager-cli/backend/internal"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func RunWails(assets embed.FS) error {
	api := internal.NewCacheManagerApi()

	return wails.Run(&options.App{
		Title:  "TabletopResourceManager",
		Width:  1024,
		Height: 1024,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        api.Startup,
		Bind:             []interface{}{api},
	})
}
