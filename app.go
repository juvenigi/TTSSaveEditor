package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GoGreet returns a greeting for the given name
func (a *App) GoGreet(name string) string {
	log.Printf("Hello %s, It's show time! (from Golang)", name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GoListDir() ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		names = append(names, dir+"/"+entry.Name())
	}
	return names, nil
}
