package main

import (
	"TTSBundler/backend/rest"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	fmt.Printf("Wonderful\n")
	if err := rest.CreateBackend().Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
