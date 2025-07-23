package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigFile = ".gatorconfig.json"
)

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir,ConfigFile)

	return path, nil
}