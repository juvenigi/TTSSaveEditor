package main

import (
	service2 "TTSBundler/service"
	"context"
)

type SaveFileAPI struct {
	ctx context.Context
}

func NewSaveFileAPI(ctx context.Context) *SaveFileAPI {
	return &SaveFileAPI{ctx: ctx}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *SaveFileAPI) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *SaveFileAPI) GoDeleteCard(savefileLocation string, deckJsonPath string) (error, string) {
	err, result := service2.DeleteCard(savefileLocation, deckJsonPath)
	return err, string(result)
}

func (a *SaveFileAPI) GoAddNewCard(patchJson service2.PartialCard, savefileLocation string, deckJsonPath string) (string, error) {
	err, result := service2.AddNewCard(patchJson, savefileLocation, deckJsonPath)
	return string(result), err
}

func (a *SaveFileAPI) GoPatchSavefile(path string, jsonp []byte) (string, error) {
	jsonBytes, err := service2.PatchSavefile(path, jsonp)
	return string(jsonBytes), err
}

func (a *SaveFileAPI) GoGetDirectory(path string) (*service2.DirectoryResponse, error) {
	return service2.GetEntries(path)
}
func (a *SaveFileAPI) GoGetSaveJson(path string) (string, error) {
	jsonBytes, err := service2.GetSaveJson(path)
	return string(jsonBytes), err
}

func (a *SaveFileAPI) GoGetEntries(path string) (*service2.DirectoryResponse, error) {
	return service2.GetEntries(path)
}

func (a *SaveFileAPI) GoGetSavefile(path string) (string, error) {
	jsonBytes, err := service2.GetSaveJson(path)
	return string(jsonBytes), err
}

func (a *SaveFileAPI) GoGetSaveResources(path string, cachePath string) (string, error) {
	res, err := service2.GetSaveResources(path, cachePath)
	return string(res), err
}
