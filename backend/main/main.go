package main

import (
	"TTSBundler/backend/rest"
	tcpSocket "TTSBundler/backend/tts-socket"
	"flag"
	"fmt"
	"github.com/pkg/browser"
	"log"
)

var editorAPI = flag.Bool("editorapi", false, "listen to external editor API events?")
var openBrowser = flag.Bool("browser", false, "open browser on start?")

func init() {
	flag.Parse()
}

func main() {
	if *editorAPI {
		go tcpSocket.ListenToTabletopApp(":39998")
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
