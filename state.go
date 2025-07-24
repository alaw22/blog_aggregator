package main

import (
	"github.com/alaw22/blog_aggregator/internal/config"
	"github.com/alaw22/blog_aggregator/internal/database"
)

type state struct {
	db *database.Queries
	conf *config.Config
}

