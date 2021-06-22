package database

import (
	hash "viv/structure/hash/manager"
	tclass "viv/structure/trie/class"
)

type Manager struct {
	table	*hash.Hash
	list	*tclass.Manager
}

func New() *Manager{
	return &Manager{
		table:    hash.New([]byte("database")),
		list:     tclass.New(),
	}
}