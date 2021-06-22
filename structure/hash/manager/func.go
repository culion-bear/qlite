package hash

import (
	"qlite/node"
	"viv/static"
)

func (m *Hash) Set(n node.Node) static.Error{
	err := m.table.Insert(n)
	if err == nil{
		m.table.Dilatation()
	}
	return err
}

func (m *Hash) SetX(n node.Node) {
	if !m.table.InsertX(n) {
		return
	}
	m.table.Dilatation()
}

func (m *Hash) DeleteList(key [][]byte) int{
	n := 0
	for _, v := range key{
		if m.table.Delete(v) == nil{
			n++
		}
	}
	if n != 0{
		s := m.table.GetNextPot(m.Size())
		if s < 16{
			s = 16
		}
		m.table.LessenWithSize(s)
	}
	return n
}

func (m *Hash) Delete(key []byte) static.Error{
	err := m.table.Delete(key)
	if err == nil{
		m.table.Lessen()
	}
	return err
}

func (m *Hash) Get(key []byte) (node.Node, static.Error){
	return m.table.Get(key)
}

func (m *Hash) Rename(oldKey, newKey []byte) static.Error{
	return m.table.Rename(oldKey, newKey)
}

func (m *Hash) RenameX(oldKey, newKey []byte) static.Error{
	err, flag := m.table.RenameX(oldKey, newKey)
	if err == nil && flag{
		m.table.Lessen()
	}
	return err
}

func (m *Hash) SelectNode() []node.Node{
	return m.table.SelectNode()
}

func (m *Hash) SelectKey() [][]byte{
	return m.table.SelectKey()
}

func (m *Hash) Size() int{
	return m.table.Size()
}

func (m *Hash) Cap() int{
	return m.table.GetCap()
}

func (m *Hash) Flush() {
	m.table.Flush()
}