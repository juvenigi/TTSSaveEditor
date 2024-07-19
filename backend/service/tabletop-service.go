package service

import (
	"TTSBundler/backend/domain"
	"errors"
	"path/filepath"
)

var tt = domain.NewTabletop()

func GetEntries(path string) (*DirectoryResponse, error) {
	if path != "" {
		err := tt.SetDirectory(path)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid path")
	}
	return &DirectoryResponse{
		Path:    path,
		Entries: GetSaves(tt),
	}, nil
}

// GetSaveJson returns the contents of a file specified by the full path string
// error occurs when the file is not json, cannot be opened, or the path is incorrect
func GetSaveJson(path string) ([]byte, error) {
	if filepath.Ext(path) != ".json" {
		return nil, errors.New("file is not a json")
	}

	if data, err := tt.GetFile(path); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func PatchSavefile(path string, jsonp []byte) ([]byte, error) {
	if path == "" {
		return nil, errors.New("invalid path")
	}
	original, err := GetSaveJson(path)
	if err != nil {
		return nil, err
	}
	return tt.PatchSaveFile(path, original, jsonp)
}

type PartialCard struct {
	LuaScript      string `json:"LuaScript"`
	LuaScriptState string `json:"LuaScriptState"`
}

func AddNewCard(data PartialCard, savePath string, deckPath string) (error, []byte) {
	return tt.AddCardsToDeck(data.LuaScript, data.LuaScriptState, savePath, deckPath)
}

func DeleteCard(savefileLocation string, path string) (error, []byte) {
	return tt.RemoveCardFromDeck(savefileLocation, path)
}
