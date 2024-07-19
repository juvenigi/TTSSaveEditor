package service

import (
	"TTSBundler/backend/domain"
	"encoding/json"
)

type DirectoryResponse struct {
	Path    string      `json:"path"`
	Entries []SaveEntry `json:"entries"`
}

type SaveEntry struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func GetSaves(tt *domain.Tabletop) []SaveEntry {
	var entries []SaveEntry
	for _, path := range tt.GetAllPaths() {
		file, err := tt.GetFile(path)
		if err != nil {
			// skip the file because it cannot be found
			continue
		}
		var data map[string]interface{}
		var saveName = ""
		err = json.Unmarshal(file, &data)
		if err == nil {
			value := data["SaveName"]
			if value != nil {
				switch value.(type) {
				case string:
					saveName = value.(string)
				}
			}
		}
		entries = append(entries, SaveEntry{Path: path, Name: saveName})
	}
	return entries
}
