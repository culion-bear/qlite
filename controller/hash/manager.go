package hash

import (
	db "qlite/database"
	class "qlite/struct"
	"qlite/trie"
)

type Manager struct{
	handle		db.Database
	list		*trie.Manager
	name		[]byte
	introduce	[]byte
}

func New(d db.Database) (trie.TrieManager, error){
	m := &Manager{
		handle:    d,
		list:      trie.New(),
		name:      []byte("db"),
		introduce: []byte("Database core structure"),
	}
	return m, m.push()
}

func (m *Manager) GetName() []byte{
	return m.name
}

func (m *Manager) GetIntroduce() []byte{
	return m.introduce
}

func (m *Manager) Work(opt []byte,l []class.Message) class.Message{
	return m.list.Work(opt, l)
}