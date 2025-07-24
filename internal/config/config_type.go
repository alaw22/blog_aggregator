package config

import "fmt"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	// Set username
	c.CurrentUser = username

	// Write data to config file
	err := write(c)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been set\n",username)


	return nil

}