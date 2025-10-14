package main

import (
	"ReallyDumbCopyPaste/tabletop_save"
	"ReallyDumbCopyPaste/util"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer util.ShowWindowOnPanic()
	gameDir, err := GetDefaultGameDir()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !strings.EqualFold(filepath.Base(wd), "PreserveMe") {
		panic(errors.New("this program must be located in a folder named PreserveMe"))
	}
	wd = filepath.Dir(wd)

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
	var cache tabletop_save.GameCacheSnooper
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

	mask := cache.GetUsedFilesMask()
	fmt.Println("mask size:", len(mask))
	cacheItems, err := DeepLs(cacheDir)
	if err != nil {
		return nil, err
	}

	var deleteList []string
	for _, entry := range cacheItems {
		base := filepath.Base(entry)
		ext := filepath.Ext(base)
		cacheFilename := tabletop_save.GetCacheFilename(strings.TrimSuffix(base, ext))
		if _, ok := mask[cacheFilename]; !ok {
			deleteList = append(deleteList, entry)
		}
	}

	return deleteList, nil
}
