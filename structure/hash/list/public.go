package list

import (
	"qlite/node"
	"viv/static"
)

//当node不存在时，加入node
func (m *Manager) Insert(v node.Node) static.Error{
	l := newNode(v)
	if m.head == nil{
		m.head, m.tail = l, l
		m.size = 1
		return nil
	}
	if m.getNode(v) != nil{
		return static.ErrNodeIsExist
	}
	m.size ++
	l.last = m.tail
	m.tail.next = l
	m.tail = l
	return nil
}

//加入node
func (m *Manager) InsertX(v node.Node) bool{
	if m.head == nil{
		l := newNode(v)
		m.head, m.tail = l, l
		m.size = 1
		return true
	}
	if n := m.getNode(v); n != nil{
		n.value = v
		return false
	}
	m.size ++
	l := newNode(v)
	l.last = m.tail
	m.tail.next = l
	m.tail = l
	return true
}

func (m *Manager) Update(key node.Node) static.Error{
	for v := m.head; v != nil; v = v.next{
		if v.value.Compare(key){
			v.value = key
			return nil
		}
	}
	return static.ErrNodeNotExist
}

func (m *Manager) SelectNode() []node.Node{
	l := make([]node.Node, m.size, m.size)
	for v, i := m.head, 0; v != nil; v, i = v.next, i + 1{
		l[i] = v.value
	}
	return l
}

func (m *Manager) SelectKey() [][]byte{
	l := make([][]byte, m.size, m.size)
	for v, i := m.head, 0; v != nil; v, i = v.next, i + 1{
		l[i] = v.value.GetKey()
	}
	return l
}

func (m *Manager) Get(key []byte) (node.Node, static.Error){
	for v := m.head; v != nil; v = v.next{
		if v.value.CompareKey(key){
			return v.value, nil
		}
	}
	return nil, static.ErrNodeNotExist
}

func (m *Manager) Delete(key []byte) static.Error{
	for v := m.head; v != nil; v = v.next{
		if v.value.CompareKey(key){
			m.deleteNode(v)
			return nil
		}
	}
	return static.ErrNodeNotExist
}

func (m *Manager) GetListNode(key []byte) (*List, node.Node, static.Error){
	for v := m.head; v != nil; v = v.next{
		if v.value.CompareKey(key){
			return v, v.value, nil
		}
	}
	return nil, nil, static.ErrNodeNotExist
}

func (m *Manager) DeleteListNode(n *List){
	m.deleteNode(n)
}

func (m *Manager) Flush() {
	m.size = 0
	m.head, m.tail = nil, nil
}