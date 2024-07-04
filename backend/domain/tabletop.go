package domain

import (
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
	saveFile []byte
}

func NewTabletop() *Tabletop {
	return &Tabletop{saves: make(map[string]*SaveFile)}
}

// getFile checks for the file in the cache and reads it from disk if not found
func (tt *Tabletop) GetFile(path string) ([]byte, error) {
	if data := tt.trySaveFromCache(path); data != nil {
		return data, nil
	}
	return tt.getSaveFromFS(path)
}

func (tt *Tabletop) trySaveFromCache(path string) []byte {
	tt.rwLock.RLock()
	cached := tt.saves[path]
	tt.rwLock.RUnlock()

	if cached != nil {
		return cached.saveFile
	}
	return nil
}

func (tt *Tabletop) getSaveFromFS(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tt.rwLock.Lock()
	tt.saves[path] = &SaveFile{saveFile: data}
	tt.rwLock.Unlock()

	return data, nil
}
