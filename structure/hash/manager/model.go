package hash

import (
	"bytes"
	"hash/crc32"
	"qlite/node"
	"viv/structure/hash/array"
)

var hashType = []byte("hash")

type Hash struct {
	key		[]byte
	table	*array.Manager
}

func New(key []byte) *Hash{
	return &Hash{
		key:   key,
		table: array.New(16),
	}
}

func (m *Hash) ToHash(size int) int{
	hashNumber:=int(crc32.ChecksumIEEE(m.key))
	if hashNumber>=0{
		return hashNumber%size
	}
	return (-hashNumber)%size
}

func (m *Hash) GetKey() []byte{
	return m.key
}

func (m *Hash) SetKey(key []byte){
	m.key = key
}

func (m *Hash) GetType() []byte{
	return hashType
}

func (m *Hash) Compare(b node.Node) bool{
	if b == nil{
		return false
	}
	return bytes.Equal(m.key, b.GetKey())
}

func (m *Hash) CompareKey(key []byte) bool{
	return bytes.Equal(m.key, key)
}