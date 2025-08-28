package properties

import (
	"os"
	"path/filepath"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/properties/internal"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
)

const PackDataYaml = resource_map.PackDataYaml

type ApplicationProperties struct {
	GameDir     string `yaml:"game-dir"`
	PackDataDir string `yaml:"pack-data-dir"`
}

type ApplicationPropertiesView struct {
	OsPathSeparator string
	GameDir         string
	PackDataDir     string
}

func (p *ApplicationProperties) InitFrom(v *ApplicationPropertiesView) error {
	p.PackDataDir = v.PackDataDir
	p.GameDir = v.GameDir

	return p.Validate()
}

func (p *ApplicationProperties) ToView() ApplicationPropertiesView {
	return ApplicationPropertiesView{
		OsPathSeparator: string(os.PathSeparator),
		GameDir:         p.GameDir,
		PackDataDir:     p.PackDataDir,
	}
}

func (p *ApplicationProperties) normalize() {
	p.PackDataDir = filepath.FromSlash(p.PackDataDir)
	p.GameDir = filepath.FromSlash(p.GameDir)
}

func (p *ApplicationProperties) Validate() error {
	if _, err := os.Lstat(p.GameDir); err != nil {
		return err
	}

	entries, err := os.ReadDir(p.PackDataDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.Compare(entry.Name(), PackDataYaml) == 0 {
			return nil
		}
	}
	packDataYamlLoc := filepath.Join(p.PackDataDir, PackDataYaml)
	if err = resource_map.CreateBlankPackDataYaml(packDataYamlLoc); err != nil {
		return err
	}
	return nil
}

func SetupDefaultAppProperties() (ApplicationProperties, error) {
	var blank ApplicationProperties
	gameDir, err := internal.GetDefaultGameDir()
	if err != nil {
		return blank, err
	}
	packData, err := internal.CreateDefaultPackData()
	if err != nil {
		return blank, err
	}

	return ApplicationProperties{
		GameDir:     gameDir,
		PackDataDir: packData,
	}, nil
}
