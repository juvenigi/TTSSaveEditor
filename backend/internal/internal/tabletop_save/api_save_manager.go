package tabletop_save

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"tts-cache-manager-cli/util"
)

// GetTabletopSaveFilesAsync todo: consider breaking packData-save-properties loop some other way, or put packData and saveProperties into the same package
func GetTabletopSaveFilesAsync(gameDir string, packDataDir string) (error, chan util.Result[TSSaveFile]) {
	resChan := make(chan util.Result[TSSaveFile])
	filenames, err := WalkDirForSaveNames(gameDir)
	if err != nil {
		close(resChan)
		return err, resChan
	}

	go func() {
		defer close(resChan)

		var wg = new(sync.WaitGroup)
		for _, fileName := range filenames {

			wg.Add(1)
			go func(fileN string) {
				defer wg.Done()
				if save, err := GetTabletopSaveFile(fileN, packDataDir); err != nil {
					resChan <- util.Result[TSSaveFile]{Err: err}
				} else {
					resChan <- util.Result[TSSaveFile]{Result: save}
				}
			}(fileName)
		}

		wg.Wait()
	}()

	return nil, resChan
}

func WalkDirForSaveNames(gameDir string) ([]string, error) {
	saveFileDir := filepath.Join(gameDir, "Saves")

	var saves []string
	root := saveFileDir + string(os.PathSeparator)
	// the death happens here, but why
	if err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		} else if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			saves = append(saves, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if len(saves) > 0 {
		if metadataIdx := slices.IndexFunc(saves, findSaveData); metadataIdx != -1 {
			saves = append(saves[:metadataIdx], saves[metadataIdx+1:]...)
		}

	}
	return saves, nil
}

func findSaveData(x string) bool {
	return strings.Compare(filepath.Base(x), "SaveFileInfos.json") == 0
}
