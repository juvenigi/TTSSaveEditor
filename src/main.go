package main

import (
	"ReallyDumbCopyPaste/tabletop_save"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	gameDir, err := GetDefaultGameDir()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = Reconcile(gameDir, wd, false)
	if err != nil {
		panic(err)
	}
	err = Reconcile(wd, gameDir, false)
	if err != nil {
		panic(err)
	}
}

func GetUnusedResources(saveFiles []string, cacheDir string) ([]string, error) {
	var cache tabletop_save.GameCacheFinder
	if err := cache.InitGameCacheFinder(cacheDir); err != nil {
		return nil, err
	}

	fmt.Println("here")

	for _, entry := range saveFiles {
		var save tabletop_save.TabletopSave
		if err := save.Init(entry); err != nil {
			return nil, err
		}
		save.Patch(&cache)
	}

	mask := cache.GetMask()
	items, err := DeepLs(cacheDir)
	if err != nil {
		return nil, err
	}

	var deleteList []string
	for _, entry := range items {
		if _, ok := mask[tabletop_save.GetCacheFilename(filepath.Base(entry))]; !ok {
			deleteList = append(deleteList, entry)
		}
	}

	return deleteList, nil
}
