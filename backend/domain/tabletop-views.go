package domain

type DirectoryResponse struct {
	Path    string      `json:"path"`
	Entries []SaveEntry `json:"entries"`
}

type SaveEntry struct {
	Path string `json:"path"`
}

func (tt *Tabletop) GetSaves() []SaveEntry {
	var entries []SaveEntry
	for path := range tt.saves {
		entries = append(entries, SaveEntry{Path: path})
	}
	return entries
}
