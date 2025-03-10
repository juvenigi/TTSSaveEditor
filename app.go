package main

import "context"

type App struct {
	ctx context.Context
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func NewApp(ctx context.Context) *App {
	return &App{ctx: ctx}
}
