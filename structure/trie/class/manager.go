package class

import (
	"errors"
	qnode "qlite/node"
)

type Manager struct {
	head	*node
}

func New() *Manager{
	return &Manager{head:newNode()}
}

func (m *Manager) Push(msg []byte, f qnode.NewNode) error{
	if len(msg) > 10 || len(msg) == 0{
		return errors.New("option is too long or too short")
	}
	return m.head.push(msg, f, 0)
}

func (m *Manager) Work(name, value []byte) qnode.Node{
	return m.head.work(name, value, 0, len(name))
}