package config

import (
	"os"
	"encoding/json"
	"fmt"
)

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Couldn't get filepath: %w\n",err)
	}

	data, err := os.ReadFile(filePath)
	
	if err != nil {
		return Config{}, fmt.Errorf("Couldn't read config file: %w\n",err)
	}

	config := Config{}

	err = json.Unmarshal(data,&config)
	if err != nil {
		return Config{}, fmt.Errorf("Couldn't Unmarshal data: %w\n",err)
	}
	return config, nil
}