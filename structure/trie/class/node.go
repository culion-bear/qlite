package class

import (
	"errors"
	qnode "qlite/node"
)

type node struct {
	list	[26]*node
	value	qnode.NewNode
}

func newNode() *node{
	return &node{}
}

func (m *node) push(msg []byte, f qnode.NewNode, k int) error{
	if k == len(msg){
		if m.value != nil{
			return errors.New(string(msg) + " structure is exist")
		}
		m.value = f
		return nil
	}
	key := msg[k] - 'a'
	if key >= 26 || key < 0{
		return errors.New(string(msg) + " structure is illegal")
	}
	if m.list[key] == nil{
		m.list[key] = newNode()
	}
	return m.list[key].push(msg, f, k + 1)
}

func (m *node) work(name, value []byte, k, length int) qnode.Node{
	if k == length{
		if m.value == nil{
			return nil
		}
		return m.value(value)
	}
	key := name[k] - 'a'
	if key >= 26 || key < 0{
		return nil
	}
	if m.list[key] == nil {
		return nil
	}
	return m.list[key].work(name, value, k + 1, length)
}