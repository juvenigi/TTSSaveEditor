package tabletop_save

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// GameCacheFinder is a map from escaped url name -> absolute fs path of cached resource
type GameCacheFinder struct {
	ResourceCached map[string]string
}

func (rv *GameCacheFinder) InitGameCacheFinder(gameDir string) error {
	rv.ResourceCached = make(map[string]string)

	if err := filepath.WalkDir(gameDir+string(filepath.Separator), func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			name := d.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))
			rv.ResourceCached[name] = path

		} else if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (rv *GameCacheFinder) Update(res *GameResource) {
	if res.Status == Local || res.Status == RemoteCached || res.Status == Packed {
		return
	}

	cacheFilename := GetCacheFilename(res.ResourceUrl)
	if _, ok := rv.ResourceCached[cacheFilename]; ok {
		res.Status = RemoteCached
	}

	return
}
