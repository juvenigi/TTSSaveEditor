package main

import (
	"TTSBundler/service"
	"context"
)

type App struct {
	ctx context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{ctx: ctx}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GoDeleteCard(savefileLocation string, deckJsonPath string) (error, string) {
	err, result := service.DeleteCard(savefileLocation, deckJsonPath)
	return err, string(result)
}

func (a *App) GoAddNewCard(patchJson service.PartialCard, savefileLocation string, deckJsonPath string) (string, error) {
	err, result := service.AddNewCard(patchJson, savefileLocation, deckJsonPath)
	return string(result), err
}

func (a *App) GoPatchSavefile(path string, jsonp []byte) (string, error) {
	jsonBytes, err := service.PatchSavefile(path, jsonp)
	return string(jsonBytes), err
}

func (a *App) GoGetDirectory(path string) (*service.DirectoryResponse, error) {
	return service.GetEntries(path)
}
func (a *App) GoGetSaveJson(path string) (string, error) {
	jsonBytes, err := service.GetSaveJson(path)
	return string(jsonBytes), err
}

func (a *App) GoGetEntries(path string) (*service.DirectoryResponse, error) {
	return service.GetEntries(path)
}

func (a *App) GoGetSavefile(path string) (string, error) {
	jsonBytes, err := service.GetSaveJson(path)
	return string(jsonBytes), err
}

func (a *App) GoGetSaveResources(path string, cachePath string) (string, error) {
	res, err := service.GetSaveResources(path, cachePath)
	return string(res), err
}
