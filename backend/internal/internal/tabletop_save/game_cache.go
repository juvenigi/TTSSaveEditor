package tabletop_save

import (
	"io/fs"
	"path/filepath"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save/internal"
)

// GameCacheFinder is a map from escaped url name -> absolute fs path of cached resource
type GameCacheFinder struct {
	ResourceCached map[string]string
}

func InitGameCacheFinder(gameDir string) (GameCacheFinder, error) {
	var result GameCacheFinder
	cachedEntries := make(map[string]string)

	err := filepath.WalkDir(gameDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			name := d.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))

			cachedEntries[name] = path
		}
		return nil
	})
	if err != nil {
		return result, err
	}

	result.ResourceCached = cachedEntries
	return result, nil
}

func (rv *GameCacheFinder) MutCacheStatus(res *GameResource) {
	if res.Status == Local || res.Status == RemoteCached || res.Status == Packed {
		return
	}

	base := filepath.Base(res.ResourceUrl)
	cacheFilename := internal.GetCacheFilename(base)
	if _, ok := rv.ResourceCached[cacheFilename]; ok {
		res.Status = RemoteCached
	}

	return
}
