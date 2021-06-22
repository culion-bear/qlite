package array

import "viv/structure/hash/list"

type Manager struct {
	used		int
	size		int
	table		[]*list.Manager
}

func New(size int) *Manager {
	return &Manager{
		used:     0,
		size:     size,
		table:    make([]*list.Manager, size, size),
	}
}