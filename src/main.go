package main

import (
	"ReallyDumbCopyPaste/tabletop_save"
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

// todo: cleanup, probably better to ask for confirmation before deleting a truckload of files
func GetUnusedResources(gameDir string) ([]string, error) {
	saveDir := filepath.Join(gameDir, "Saves")
	workshopDir := filepath.Join(gameDir, "Mods", "Workshop")

	entries, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, err
	}
	workshopEntries, err := os.ReadDir(workshopDir)
	if err != nil {
		return nil, err
	}
	entries = append(entries, workshopEntries...)

	var cache tabletop_save.GameCacheFinder
	if err := cache.InitGameCacheFinder(gameDir); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		var save tabletop_save.TabletopSave
		if err := save.Init(entry.Name()); err != nil {
			return nil, err
		}
		save.Patch(&cache)
	}

	mask := cache.GetMask()

	// todo: exclude Workshop dir
	items, err := DeepLs(filepath.Join(gameDir, "Mods"))
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
