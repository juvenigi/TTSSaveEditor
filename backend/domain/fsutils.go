package domain

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func parseEntries(entries []os.DirEntry, currentDir string, visited map[string]bool, saveFiles map[string]*SaveFile) []string {
	var newQueue []string
	for _, entry := range entries {
		absPath, err := filepath.Abs(currentDir)
		if err != nil {
			continue
		}
		fullPath := filepath.Join(absPath, entry.Name())
		log.Infof("entry: %s", fullPath)
		if entry.IsDir() {
			if !visited[fullPath] {
				visited[fullPath] = true
				newQueue = append(newQueue, fullPath)
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		file, err := os.Open(fullPath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", fullPath, err)
			continue
		}
		_ = file.Close()
		saveFiles[fullPath] = &SaveFile{stale: false}
	}
	return newQueue
}

// TODO makes more sense to keep here, but I am using this function directly in the controller, which is not smart
func GetDefaultTTSAbsPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(usr.HomeDir, "Documents\\My Games\\Tabletop Simulator\\Saves")
	default:
		return "", errors.New("not implemented")
	}
	return path, nil
}
