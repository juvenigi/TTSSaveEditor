package properties

import (
	"bufio"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

const (
	packDataFolder = "PackData"
	PackDataYaml   = "resource-map.yaml"

	ConfigFileYaml = "properties.yaml"
)

type ApplicationProperties struct {
	GameDir     string `yaml:"game-dir"`
	PackDataDir string `yaml:"pack-data-dir"`
}

func newProperties(fileContent io.Reader) (ApplicationProperties, error) {
	properties := ApplicationProperties{}

	if err := yaml.NewDecoder(bufio.NewReader(fileContent)).Decode(&properties); err != nil {
		return properties, err
	}
	if err := properties.validate(); err != nil {
		return properties, err
	}

	return properties, nil
}

func (p *ApplicationProperties) validate() error {
	if _, err := os.Lstat(p.GameDir); err != nil {
		return err
	}

	dir, packDataIndex := filepath.Split(p.PackDataDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.Compare(entry.Name(), packDataIndex) == 0 {
			return nil
		}
	}
	return errors.New("could not find pack data")
}

func getDefaultGameDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Join(errors.New("unable to get home dir"), err)
	}

	switch runtime.GOOS {
	case "windows", "darwin":
		return filepath.Join(home, "Documents", "My games", "Tabletop Simulator"), nil
	case "linux":
		return filepath.Join(home, ".local", "share", "Tabletop Simulator"), nil
	default:
		return "", errors.New("unknown default save data location for " + runtime.GOOS)
	}
}

func createDefaultPackData() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", errors.Join(errors.New("unable to get executable path"), err)
	}
	resPackPath := filepath.Join(filepath.Dir(exePath), packDataFolder)
	if err = os.Mkdir(resPackPath, 0755); err != nil && !os.IsExist(err) {
		return "", errors.Join(errors.New("unable to create resource data dir"), err)
	}
	resourcePack := filepath.Join(resPackPath, PackDataYaml)
	create, err := os.Create(resourcePack)
	if err != nil {
		panic(err)
	}
	err = create.Close()
	if err != nil {
		return "", err
	}
	return resourcePack, nil
}

func writeDefaultProperties() (string, error) {
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

	configFilepath := filepath.Join(filepath.Dir(executable), ConfigFileYaml)
	err = os.WriteFile(configFilepath, blob, 0755)
	if err != nil {
		return "", err
	}

	return configFilepath, nil

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
	if slices.Contains(files, ConfigFileYaml) {
		return filepath.Join(configDir, ConfigFileYaml)
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
