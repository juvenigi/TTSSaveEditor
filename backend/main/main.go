package main

import (
	"TTSBundler/backend/rest"
	tcpSocket "TTSBundler/backend/tts-socket"
	"log"
)

func main() {
	go tcpSocket.ListenToTabletopApp(":39998")

	if err := rest.CreateSaveEditorBackend().Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
