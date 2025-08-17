package internal

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

const (
	packDataFolder = "PackData"

	ConfigFileYaml = "properties.yaml"
)

func GetDefaultGameDir() (string, error) {
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

func CreateDefaultPackData() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", errors.Join(errors.New("unable to get executable path"), err)
	}
	resPackPath := filepath.Join(filepath.Dir(exePath), packDataFolder)
	if err = os.Mkdir(resPackPath, 0755); err != nil && !os.IsExist(err) {
		return "", errors.Join(errors.New("unable to create resource data dir"), err)
	}
	if err != nil {
		return "", err
	}
	return resPackPath, nil
}
