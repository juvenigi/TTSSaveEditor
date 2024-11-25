package main

import (
	"TTSBundler/backend/grdive"
	"TTSBundler/backend/rest"
	tabletopSocket "TTSBundler/backend/tts-socket"
	"flag"
	"fmt"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"log"
)

var (
	begForToken = flag.Bool("gdrive", false, "open auth consent window to obtain token?")
	editorAPI   = flag.Bool("editorapi", false, "listen to external editor API events?")
	openBrowser = flag.Bool("browser", false, "open browser on start?")
)

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
		go tabletopSocket.ListenToAppTCP(":39998")
	}
	if *openBrowser {
		err := browser.OpenURL("http://localhost:3000")
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
	if err := rest.CreateSaveEditorBackend().Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
