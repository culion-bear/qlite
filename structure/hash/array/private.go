package array

import (
	"hash/crc32"
	"viv/structure/hash/list"
)

func (m *Manager) toHash(buf []byte) int{
	hashNumber := int(crc32.ChecksumIEEE(buf))
	if hashNumber >= 0{
		return hashNumber % m.size
	}
	return (-hashNumber) % m.size
}

func (m *Manager) getLoadFactor() float32{
	return float32(m.used) / float32(m.size)
}

func (m *Manager) rehash(size int) {
	l := m.SelectNode()
	m.used = 0
	m.size = size
	m.table = make([]*list.Manager, size, size)
	for _, v := range l{
		m.Insert(v)
	}
}