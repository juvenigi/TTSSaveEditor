package tabletop_save

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tts-cache-manager-cli/util"
)

func GetTabletopSaveFilesAsync(scanDir string, packDataDir string) (error, chan util.Result[TSSaveFile]) {
	resChan := make(chan util.Result[TSSaveFile])
	sameAsPackData := strings.EqualFold(packDataDir, scanDir)

	filenames, err := WalkDirForSaveNames(scanDir, sameAsPackData)
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
					resChan <- util.Result[TSSaveFile]{Result: *save}
				}
			}(fileName)
		}

		wg.Wait()
	}()

	return nil, resChan
}

// todo: fix cyclic import of pack data json
func WalkDirForSaveNames(scanDir string, isPackData bool) ([]string, error) {
	var saves []string
	metadataSkipped := false
	metadataName := "SaveFileInfos.json"
	if isPackData {
		metadataName = "resource-map.json"
	}

	walkDirFunc := func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		} else if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			if metadataSkipped || strings.Compare(info.Name(), metadataName) != 0 {
				saves = append(saves, path)
			} else {
				metadataSkipped = true
			}
		}
		return nil
	}

	saveFileDir := filepath.Join(scanDir)
	if err := filepath.WalkDir(saveFileDir+string(os.PathSeparator), walkDirFunc); err != nil {
		return nil, err
	} else {
		return saves, nil
	}
}
