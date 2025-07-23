package config

import (
	"os"
	"fmt"
)

func write(data []byte) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Couldn't get filepath %w\n",err)
	}

	err = os.WriteFile(filePath, data, 644)
	
	if err != nil {
		return fmt.Errorf("Couldn't write config file %w\n",err)
	}

	return nil
}