package database

import "qlite/node"

func (m *Manager) Push(name []byte, createFunc node.NewNode) error{
	return m.list.Push(name, createFunc)
}