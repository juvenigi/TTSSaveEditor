package service

import (
	"TTSBundler/backend/domain"
	"errors"
	"path/filepath"
)

var tt = domain.NewTabletop()

func GetEntries(path string) (*domain.DirectoryResponse, error) {
	if path != "" {
		err := tt.SetDirectory(path)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid path")
	}
	return &domain.DirectoryResponse{
		Path:    path,
		Entries: tt.GetSaves(),
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
