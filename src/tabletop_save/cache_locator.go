package tabletop_save

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var excludeChars = regexp.MustCompile(`[^a-zA-Z0-9]`)

// GetCacheFilename note: this trims file extension
// TTS decides to turn `file.png` into `filepng.png` when caching, therefore we need to account this in our code.
// this is why I chose to remove the file extension completely if it is present
// note2: sometimes people put dropbox links
// like this: `https://www.dropbox.com/s/<id>/image.png?dl=1`
// tabletop simulator creates `httpswwwdropboxcomsidimagepngdl1.png`
// this breaks a lot of previously held assumptions, meaning I do need to handle the entire input string
func GetCacheFilename(raw string) string {
	return excludeChars.ReplaceAllString(raw, "")
}

// GameCacheSnooper is a map from escaped url name -> absolute fs path of cached resource
type GameCacheSnooper struct {
	ResourceUsed map[string]bool
}

func (rv *GameCacheSnooper) InitGameCacheFinder(gameDir string) error {
	rv.ResourceUsed = make(map[string]bool)

	if err := filepath.WalkDir(gameDir+string(filepath.Separator), func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			name := d.Name()
			// file extensions exist depending on a vibe,
			// so let's always discard them for the purpose of finding suitable files
			rv.ResourceUsed[strings.TrimSuffix(name, filepath.Ext(name))] = false

		} else if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (rv *GameCacheSnooper) MarkAsUsed(resUrl string) {
	cacheFilename := GetCacheFilename(resUrl)
	rv.ResourceUsed[cacheFilename] = true
}

func (rv *GameCacheSnooper) GetUsedFilesMask() map[string]struct{} {
	res := make(map[string]struct{})
	for k, resourceUsed := range rv.ResourceUsed {
		if resourceUsed {
			res[k] = struct{}{}
		}
	}

	return res
}

func GetAllSaveLocations(defaultGameDir string) []string {
	var result []string
	var err error
	saveDir := filepath.Join(defaultGameDir, "Saves")
	workshopDir := filepath.Join(defaultGameDir, "Mods", "Workshop")

	result, err = getSaveLocations(saveDir)
	if err != nil {
		return nil
	}

	workshopLocs, err := getSaveLocations(workshopDir)
	if err == nil {
		result = append(result, workshopLocs...)
	}

	return result
}

func getSaveLocations(dir string) ([]string, error) {
	var saveFiles []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(entryName, ".json") || strings.EqualFold(entryName, "SaveFileInfos.json") {
			continue
		}
		saveFiles = append(saveFiles, filepath.Join(dir, entryName))
	}

	return saveFiles, nil
}
