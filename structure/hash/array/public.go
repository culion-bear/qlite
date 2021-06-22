package array

import (
	"qlite/node"
	"viv/static"
	"viv/structure/hash/list"
)

func (m *Manager) Size() int{
	return m.used
}

func (m *Manager) SelectNode() []node.Node{
	l := make([]node.Node, 0)
	for _, v := range m.table{
		if v == nil{
			continue
		}
		l = append(l, v.SelectNode()...)
	}
	return l
}

func (m *Manager) SelectKey() [][]byte{
	l := make([][]byte, 0)
	for _, v := range m.table{
		if v == nil{
			continue
		}
		l = append(l, v.SelectKey()...)
	}
	return l
}

func (m *Manager) Insert(v node.Node) static.Error{
	h := v.ToHash(m.size)
	if m.table[h] == nil{
		m.table[h] = list.New()
	}
	err := m.table[h].Insert(v)
	if err != nil{
		return err
	}
	m.used ++
	return nil
}

func (m *Manager) InsertX(v node.Node) bool{
	h := v.ToHash(m.size)
	if m.table[h] == nil{
		m.table[h] = list.New()
	}
	if m.table[h].InsertX(v){
		m.used ++
		return true
	}
	return false
}

func (m *Manager) Update(v node.Node) static.Error{
	h := v.ToHash(m.size)
	if m.table[h] == nil{
		return static.ErrNodeNotExist
	}
	return m.table[h].Update(v)
}

func (m *Manager) Delete(key []byte) static.Error{
	h := m.toHash(key)
	if m.table[h] == nil{
		return static.ErrNodeNotExist
	}
	err := m.table[h].Delete(key)
	if err == nil{
		m.used --
	}
	return err
}

func (m *Manager) Get(key []byte) (node.Node, static.Error){
	h := m.toHash(key)
	if m.table[h] == nil{
		return nil, static.ErrNodeNotExist
	}
	return m.table[h].Get(key)
}

//重命名，仅当newKey不存在
func (m *Manager) Rename(oldKey, newKey []byte) static.Error{
	oldHash := m.toHash(oldKey)
	if m.table[oldHash] == nil{
		return static.ErrNodeNotExist
	}
	l, n, err := m.table[oldHash].GetListNode(oldKey)
	if err != nil{
		return err
	}
	n.SetKey(newKey)
	newHash := m.toHash(newKey)
	if m.table[newHash] == nil{
		m.table[newHash] = list.New()
	}
	err = m.table[newHash].Insert(n)
	if err != nil{
		n.SetKey(oldKey)
		return err
	}
	m.table[oldHash].DeleteListNode(l)
	return nil
}

//重命名
func (m *Manager) RenameX(oldKey, newKey []byte) (static.Error, bool){
	oldHash := m.toHash(oldKey)
	if m.table[oldHash] == nil{
		return static.ErrNodeNotExist, false
	}
	l, n, err := m.table[oldHash].GetListNode(oldKey)
	if err != nil{
		return err, false
	}
	n.SetKey(newKey)
	newHash := m.toHash(newKey)
	if m.table[newHash] == nil{
		m.table[newHash] = list.New()
	}
	if !m.table[newHash].InsertX(n){
		m.used --
		return nil, true
	}
	m.table[oldHash].DeleteListNode(l)
	return nil, false
}

func (m *Manager) Flush() {
	m.used = 0
	m.size = 16
	m.table = make([]*list.Manager, 16, 16)
}