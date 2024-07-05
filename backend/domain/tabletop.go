package domain

import (
	"fmt"
	"os"
	"sync"
)

// tabletop holds a map of save files with a mutex for concurrent access
type Tabletop struct {
	saves  map[string]*SaveFile
	rwLock sync.RWMutex
}

// SaveFile holds the byte data of a save file
type SaveFile struct {
	file     *os.File
	saveData []byte
	stale    bool
}

func NewTabletop() *Tabletop {
	return &Tabletop{saves: make(map[string]*SaveFile)}
}

func NewSavefile(fullPath string, data []byte) (*SaveFile, error) {
	err := os.WriteFile(fullPath, data, 0666)
	if err != nil {
		return nil, err
	}
	return &SaveFile{saveData: data}, err
}

// ScanPathForSaves scans a directory and its subdirectories for JSON files,
// returning a map where the key is the full file path and the value is the SaveFile struct.
func (tt *Tabletop) ScanPathForSaves(path string) error {
	saveFiles := make(map[string]*SaveFile)
	visited := make(map[string]bool)
	// Queue for directories to process
	queue := []string{path}
	for len(queue) > 0 {
		currentDir := queue[0]
		queue = queue[1:]

		if entries, err := os.ReadDir(currentDir); err != nil {
			fmt.Printf("Error reading directory %s: %v\n", currentDir, err)
			continue
		} else {
			queue = append(queue, parseEntries(entries, currentDir, visited, saveFiles)...)
		}
	}
	tt.saves = saveFiles
	return nil
}

// GetFile checks for the file in the cache and reads it from disk if not found
func (tt *Tabletop) GetFile(path string) ([]byte, error) {
	if data := tt.trySaveFromCache(path); data != nil {
		return data, nil
	}
	return tt.getSaveFromFS(path)
}

func (tt *Tabletop) trySaveFromCache(path string) []byte {
	tt.rwLock.RLock()
	cached, exists := tt.saves[path]
	tt.rwLock.RUnlock()

	if exists {
		return cached.saveData
	}
	return nil
}

func (tt *Tabletop) getSaveFromFS(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tt.rwLock.Lock()
	tt.saves[path] = &SaveFile{saveData: data}
	tt.rwLock.Unlock()

	return data, nil
}
