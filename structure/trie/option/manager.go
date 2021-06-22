package option

import (
	"errors"
	class "qlite/struct"
	"qlite/trie"
)

type Manager struct {
	head	*node
}

func New() *Manager{
	return &Manager{head:newNode()}
}

func (m *Manager) Push(f trie.TrieManager) error{
	if len(f.GetName()) > 10 || len(f.GetName()) == 0{
		return errors.New("option is too long or too short")
	}
	return m.head.push(f.GetName(), f, 0)
}

func (m *Manager) Work(name, opt []byte, list []class.Message) class.Message{
	return m.head.work(name, opt, list, 0, len(name))
}