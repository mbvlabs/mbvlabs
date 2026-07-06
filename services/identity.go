package services

import (
	"mbvlabs/config"
	"mbvlabs/internal/storage"
	"mbvlabs/queue"
)

type Identity struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
	pepper     string
}

func NewIdentity(db storage.Pool, insertOnly queue.InsertOnly, cfg config.Config) Identity {
	return Identity{db: db, insertOnly: insertOnly, pepper: cfg.Auth.Pepper}
}
