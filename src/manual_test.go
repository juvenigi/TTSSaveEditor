package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManually(t *testing.T) {
	src := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator\\Mods"
	dst := "C:\\Users\\joghourt\\home\\dest"

	err := Reconcile(src, dst, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeepLs(t *testing.T) {
	src := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator\\Mods"
	entries, err := DeepLs(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		t.Log(entry)
	}
}

func TestGetUnusedResources(t *testing.T) {
	var err error
	//cacheDir := "C:\\Users\\joghourt\\home\\dest"
	cacheDir := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator\\Mods"

	saveDir := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator\\Saves"

	var entries []os.DirEntry
	entries, err = os.ReadDir(saveDir)
	if err != nil {
		t.Fatal(err)
	}
	var saveFiles []string
	for _, entry := range entries {
		entryName := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(entryName, ".json") || strings.EqualFold(entryName, "SaveFileInfos.json") {
			continue
		}
		saveFiles = append(saveFiles, filepath.Join(saveDir, entryName))
	}

	resources, err := GetUnusedResources(saveFiles, cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range resources {
		t.Log(entry)
	}

	for _, entry := range resources {
		err = os.Remove(filepath.Join(cacheDir, entry))
		if err != nil {
			t.Error(err)
		}
	}

}
