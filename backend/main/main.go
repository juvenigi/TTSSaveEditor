package main

import (
	"TTSBundler/backend/rest"
	"log"
)

func main() {
	if err := rest.CreateBackend().Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
