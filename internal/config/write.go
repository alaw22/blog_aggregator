package config

import (
	"os"
	"fmt"
	"encoding/json"
)

func write(c *Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Couldn't get filepath: %w\n",err)
	}

	// u = rw, g = r, o = r -> 644
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't open config file to write: %w\n",err)
	}

	defer f.Close()

	data, err := json.MarshalIndent(c,"","    ")
	if err != nil {
		return fmt.Errorf("Error in json.MarshalIndent: %w\n",err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("Some of or none of the data was written: %w\n",err)
	}
	
	if err != nil {
		return fmt.Errorf("Couldn't write config file: %w\n",err)
	}

	return nil
}