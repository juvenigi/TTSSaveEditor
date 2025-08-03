package tabletop_save

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"tts-cache-manager-cli/properties"
	"tts-cache-manager-cli/util"
)

type SaveManager struct {
	properties properties.ApplicationProperties
}

func NewSaveManager(properties properties.ApplicationProperties) SaveManager {
	return SaveManager{properties: properties}
}

func (sm *SaveManager) GetTabletopSaveFilesAsync(filesCh chan<- TSSaveFile, errCh chan<- error) (chan struct{}, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	files, err := sm.lsDirForNames()
	if err != nil {
		return doneCh, err
	}

	go func() {
		var wg = new(sync.WaitGroup)
		for _, file := range files {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				defer util.DeadGopherChannel(errCh)

				if save, err := sm.GetTabletopSaveFile(file); err != nil {
					errCh <- err
					return
				} else {
					filesCh <- save
				}
			}(file)
		}
		wg.Wait()
	}()

	return doneCh, nil
}

func (sm *SaveManager) GetTabletopSaveFile(loc string) (TSSaveFile, error) {
	var dummy TSSaveFile

	file, err := os.ReadFile(loc)
	if err != nil {
		return dummy, err
	}

	bundles, err := ParseResourcesBundles(file, sm.properties.PackDataDir)
	if err != nil {
		return dummy, err
	}

	return TSSaveFile{
		savefileLocation: loc,
		resourceBundle:   bundles,
	}, nil
}

func (sm *SaveManager) lsDirForNames() ([]string, error) {
	saveFileDir := filepath.Join(sm.properties.GameDir, "Saves")
	dir, err := os.ReadDir(saveFileDir)
	if err != nil {
		return nil, err
	}

	var saves []string
	for _, file := range dir {
		name := file.Name()
		if strings.HasSuffix(name, ".json") {
			saves = append(saves, filepath.Join(saveFileDir, name))
		}
	}

	if metadataIdx := slices.IndexFunc(saves, findSaveData); metadataIdx != -1 {
		saves = append(saves[:metadataIdx], saves[metadataIdx+1:]...)
	}

	return saves, nil
}

func findSaveData(x string) bool {
	return strings.Compare(filepath.Base(x), "SaveFileInfos.json") == 0
}
