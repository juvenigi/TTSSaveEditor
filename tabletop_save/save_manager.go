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

func (sm *SaveManager) GetTabletopSaveFilesAsync() (error, chan util.Result[TSSaveFile]) {
	resChan := make(chan util.Result[TSSaveFile])

	files, err := sm.walkDirForSaveNames()
	if err != nil {
		close(resChan)
		return err, resChan
	}

	go func() {
		defer close(resChan)

		var wg = new(sync.WaitGroup)
		for _, file := range files {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				if save, err := sm.GetTabletopSaveFile(file); err != nil {
					resChan <- util.Result[TSSaveFile]{Err: err}
				} else {
					resChan <- util.Result[TSSaveFile]{Result: save}
				}
			}(file)
		}

		wg.Wait()
	}()

	return nil, resChan
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

func (sm *SaveManager) walkDirForSaveNames() ([]string, error) {
	saveFileDir := filepath.Join(sm.properties.GameDir, "Saves")

	var saves []string
	err := filepath.Walk(saveFileDir+string(os.PathSeparator), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			saves = append(saves, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if metadataIdx := slices.IndexFunc(saves, findSaveData); metadataIdx != -1 {
		saves = append(saves[:metadataIdx], saves[metadataIdx+1:]...)
	}

	return saves, nil
}

func findSaveData(x string) bool {
	return strings.Compare(filepath.Base(x), "SaveFileInfos.json") == 0
}
