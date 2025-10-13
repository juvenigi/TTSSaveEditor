package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func GetDefaultGameDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Join(errors.New("unable to get home dir"), err)
	}

	switch runtime.GOOS {
	case "windows", "darwin":
		return filepath.Join(home, "Documents", "My games", "Tabletop Simulator", "Mods"), nil
	case "linux":
		return filepath.Join(home, ".local", "share", "Tabletop Simulator", "Mods"), nil
	default:
		return "", errors.New("unknown default save data location for " + runtime.GOOS)
	}
}
