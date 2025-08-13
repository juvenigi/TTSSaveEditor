package properties

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/properties/internal"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"

	"gopkg.in/yaml.v3"
)

const PackDataYaml = internal.PackDataYaml

type ApplicationProperties struct {
	GameDir     string `yaml:"game-dir"`
	PackDataDir string `yaml:"pack-data-dir"`
}

type ApplicationPropertiesView struct {
	OsPathSeparator string
	GameDir         string
	PackDataDir     string
}

func (p *ApplicationProperties) ToView() ApplicationPropertiesView {
	return ApplicationPropertiesView{
		OsPathSeparator: string(os.PathSeparator),
		GameDir:         p.GameDir,
		PackDataDir:     p.PackDataDir,
	}
}

func newProperties(fileContent io.Reader) (ApplicationProperties, error) {
	properties := ApplicationProperties{}

	if err := yaml.NewDecoder(bufio.NewReader(fileContent)).Decode(&properties); err != nil {
		return properties, err
	}

	properties.normalize()

	if err := properties.validate(); err != nil {
		return properties, err
	}

	return properties, nil
}

func (p *ApplicationProperties) normalize() {
	p.PackDataDir = filepath.FromSlash(p.PackDataDir)
	p.GameDir = filepath.FromSlash(p.GameDir)
}

func (p *ApplicationProperties) validate() error {
	if _, err := os.Lstat(p.GameDir); err != nil {
		return err
	}

	entries, err := os.ReadDir(p.PackDataDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.Compare(entry.Name(), internal.PackDataYaml) == 0 {
			return nil
		}
	}
	packDataYamlLoc := filepath.Join(p.PackDataDir, internal.PackDataYaml)
	if err = resource_map.CreateBlankPackDataYaml(packDataYamlLoc); err != nil {
		return err
	}
	return nil
}

// todo: handle this more gracefully with a flag, don't use raw position-based arg we are not in the 20th century
// returns an empty string if config file cannot be found
func getCfgFilePath() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path")
		return ""
	}
	configDir := filepath.Dir(execPath)

	files := getFilenames(configDir)
	if len(files) == 0 {
		log.Println("No config files found in ", configDir)
		return ""
	}
	if slices.Contains(files, internal.ConfigFileYaml) {
		return filepath.Join(configDir, internal.ConfigFileYaml)
	}

	return ""
}

func getFilenames(configDir string) []string {
	files, err := os.ReadDir(configDir)
	if err != nil {
		log.Printf("error reading config directory %s: %v", configDir, err)
		return nil
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames
}

func GetDefaultApplicationProperties() (ApplicationProperties, error) {
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

func WriteDefaultProperties() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	properties, err := GetDefaultApplicationProperties()
	if err != nil {
		return "", err
	}

	blob, err := yaml.Marshal(properties)
	if err != nil {
		return "", err
	}

	configFilepath := filepath.Join(filepath.Dir(executable), internal.ConfigFileYaml)
	err = os.WriteFile(configFilepath, blob, 0755)
	if err != nil {
		return "", err
	}

	return configFilepath, nil

}
