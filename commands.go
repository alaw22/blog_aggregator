package main

import (
	"fmt"
	// "errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("Unrecognized command '%s'\n",cmd.name)
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Potential to use mutex?
func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

