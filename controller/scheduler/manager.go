package sd

import (
	"viv/controller/logistics"
	"viv/structure/trie/option"
)

var Scheduler *Manager

type Manager struct {
	password	[]byte
	ch			chan *Model
	tree		*option.Manager
	log			*logistics.Manager
}

func New(password []byte) *Manager{
	return &Manager{
		password: password,
		ch:       make(chan *Model, 100000),
		tree:     option.New(),
	}
}