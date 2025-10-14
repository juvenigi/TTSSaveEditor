package main

import (
	"ReallyDumbCopyPaste/tabletop_save"
	"ReallyDumbCopyPaste/util"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer util.ShowWindowOnPanic()

	dummy := bufio.NewReader(os.Stdin)

	gameDir, err := GetDefaultGameDir()
	if err != nil {
		panic(err)
	}

	savefiles := tabletop_save.GetAllSaveLocations(gameDir)

	unused, err := GetUnusedResources(savefiles, gameDir)
	if err != nil {
		panic(err)
	}
	fmt.Println("Unused resources:")
	if len(unused) == 0 {
		fmt.Println("No unused resources")
	} else {
		for _, v := range unused {
			fmt.Println(v)
		}
		fmt.Println("press enter to continue")
		_, _ = dummy.ReadString('\n')
	}
	fmt.Println("Deleting...")
	for _, entry := range unused {
		err = os.Remove(filepath.Join(gameDir, entry))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("press enter to exit")
	_, _ = dummy.ReadString('\n')
}

func GetUnusedResources(saveFiles []string, cacheDir string) ([]string, error) {
	var cache tabletop_save.GameCacheSnooper
	if err := cache.InitGameCacheFinder(cacheDir); err != nil {
		return nil, err
	}

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
