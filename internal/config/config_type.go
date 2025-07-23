package config

import (
	"encoding/json"
)

type Config struct {
	DB_URL string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	// Set username
	c.CurrentUser = username

	// Marshal data to bytes
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Write data to config file
	err = write(data)
	if err != nil {
		return err
	}

	return nil

}