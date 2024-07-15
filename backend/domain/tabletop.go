package domain

import (
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// tabletop holds a map of save files with a mutex for concurrent access
type Tabletop struct {
	saves  map[string]*SaveFile
	rwLock sync.RWMutex
}

// SaveFile holds the byte data of a save file
type SaveFile struct {
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

// SetDirectory scans a directory and its subdirectories for JSON files,
// returning a map where the key is the full file path and the value is the SaveFile struct.
func (tt *Tabletop) SetDirectory(path string) error {
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

func (tt *Tabletop) PatchSaveFile(path string, original []byte, patches []byte) ([]byte, error) {
	patch, err := jsonpatch.DecodePatch(patches)
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

func (tt *Tabletop) BackupSaveFile(path string, data []byte) error {
	save := tt.saves[path]
	if save == nil {
		return errors.New(path + " could not be found")
	}
	// backup previous file
	dir, name := filepath.Split(path)
	baknum := getLargestBakNum(dir)
	err := os.WriteFile(filepath.Join(dir, fmt.Sprint(name, ".", baknum, ".bak")), save.saveData, 0666)
	if err != nil {
		return err
	}
	// write new file
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	// update cache
	save.saveData = data
	save.stale = true
	return nil
}

func getLargestBakNum(dir string) int {
	largestBak := 0
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		segments := strings.Split(entry.Name(), ".")
		slen := len(segments)
		if slen > 4 && slen-2 > 0 {
			bakNumStr := segments[slen-2]
			num, err := strconv.ParseInt(bakNumStr, 10, 8)
			if err != nil {
				continue
			}
			largestBak = int(num)
		}
	}
	return largestBak
}
