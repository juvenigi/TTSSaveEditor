package tabletop_save

import (
	"io/fs"
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

// GameCacheFinder is a map from escaped url name -> absolute fs path of cached resource
type GameCacheFinder struct {
	ResourceCached map[string]bool
}

func (rv *GameCacheFinder) InitGameCacheFinder(gameDir string) error {
	rv.ResourceCached = make(map[string]bool)

	if err := filepath.WalkDir(gameDir+string(filepath.Separator), func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			name := d.Name()
			// file extensions exist depending on a vibe,
			// so let's always discard them for the purpose of finding suitable files
			rv.ResourceCached[strings.TrimSuffix(name, filepath.Ext(name))] = false

		} else if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (rv *GameCacheFinder) MarkAsUsed(resUrl string) {
	cacheFilename := GetCacheFilename(resUrl)
	_, ok := rv.ResourceCached[cacheFilename]
	if ok {
		rv.ResourceCached[cacheFilename] = true
	}
}

func (rv *GameCacheFinder) GetMask() map[string]struct{} {
	res := make(map[string]struct{})
	for k, v := range rv.ResourceCached {
		if v {
			res[k] = struct{}{}
		}
	}

	return res
}
