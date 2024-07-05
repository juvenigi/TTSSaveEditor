package service

import (
	"TTSBundler/backend/domain"
	"encoding/json"
	"errors"
	"path/filepath"
)

var tt = domain.NewTabletop()

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

func GetEntries(path string) ([]byte, error) {
	if path != "" {
		err := tt.ScanPathForSaves(path)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(tt.GetSaves())
}
