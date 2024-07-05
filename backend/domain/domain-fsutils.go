package domain

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"path/filepath"
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
		saveFiles[fullPath] = &SaveFile{file: file, stale: true}
	}
	return newQueue
}
