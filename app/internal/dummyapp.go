package app

import (
	"context"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewDummyApp creates a new App application struct
func NewDummyApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	log.Println("Hello " + name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
