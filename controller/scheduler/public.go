package sd

import (
	"qlite/trie"
	"viv/controller/logistics"
	"viv/structure/qt"
)

func (m *Manager) GetPassword() []byte{
	return m.password
}

func (m *Manager) Push(node trie.TrieManager) error{
	return m.tree.Push(node)
}

func (m *Manager) Run(log *logistics.Manager){
	m.log = log
	go m.log.Run()
	var n *Model
	for true{
		n = <- m.ch
		n.ch <- m.work(n.msg)
	}
}

func (m *Manager) Option(n *Model){
	m.ch <- n
}

func (m *Manager) GetChan() chan *Model{
	return m.ch
}

func (m *Manager) Load(msg *qt.Queue){
	m.load(msg)
}

func (m *Manager) Close(){
	m.log.Close()
}