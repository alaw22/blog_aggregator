package main

import (
	"github.com/alaw22/gator/internal/config"
	"github.com/alaw22/gator/internal/database"
)

type state struct {
	db *database.Queries
	conf *config.Config
}

