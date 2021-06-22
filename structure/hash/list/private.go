package list

import (
	"qlite/node"
)

func (m *Manager) getNode(k node.Node) *List {
	for v := m.head; v != nil; v = v.next{
		if v.value.Compare(k) {
			return v
		}
	}
	return nil
}

func (m *Manager) deleteNode(k *List){
	m.size --
	if m.head == k && m.tail == k{
		m.head, m.tail = nil, nil
	}else if m.head == k{
		m.head = k.next
		m.head.last = nil
	}else if m.tail == k{
		m.tail = k.last
		m.tail.next = nil
	}else{
		k.last.next = k.next
		k.next.last = k.last
	}
}