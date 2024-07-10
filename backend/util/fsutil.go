package util

import (
	"errors"
	"os/user"
	"path/filepath"
	"runtime"
)

func GetDefaultTTSAbsPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(usr.HomeDir, "Documents\\My Games\\Tabletop Simulator\\Saves")
	default:
		return "", errors.New("not implemented")
	}
	return path, nil
}
