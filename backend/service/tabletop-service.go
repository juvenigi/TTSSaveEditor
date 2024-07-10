package service

import (
	"TTSBundler/backend/domain"
	"encoding/json"
	"errors"
	jsonpatch "github.com/evanphx/json-patch"
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
	} else {
		return nil, errors.New("invalid path")
	}
	return json.Marshal(domain.DirectoryResponse{
		Path:    path,
		Entries: tt.GetSaves(),
	})
}

func PatchSavefile(path string, jsonp []byte) ([]byte, error) {
	if path == "" {
		return nil, errors.New("invalid path")
	}
	original, err := GetSaveJson(path)
	if err != nil {
		return nil, err
	}
	patch, err := jsonpatch.DecodePatch(jsonp)
	if err != nil {
		return nil, err
	}
	modified, err := patch.Apply(original)
	if err != nil {
		return nil, err
	}
	err = tt.BackupSaveFile(path, modified)
	if err != nil {
		return nil, err
	}

	return modified, nil
}
