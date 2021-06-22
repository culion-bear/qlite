package qt

import class "qlite/struct"

type Tree struct {
	list	*Queue
	value	[]class.Message
}

func NewTree() *Tree{
	return &Tree{
		list:  NewQueue(),
		value: make([]class.Message, 0),
	}
}

func (m *Tree) PushValue(msg class.Message){
	m.value = append(m.value, msg)
}

func (m *Tree) PushNode(msg *Tree){
	m.list.Push(msg, len(m.value))
	m.value = append(m.value, class.Message{})
}

func (m *Tree) GetValue() []class.Message{
	return m.value
}

func (m *Tree) GetNode() (*Tree, int){
	return m.list.Pop()
}

func (m *Tree) SetValue(msg class.Message, p int){
	m.value[p] = msg
}