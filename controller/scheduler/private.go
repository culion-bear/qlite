package sd

import (
	class "qlite/struct"
	"viv/static"
	"viv/structure/qt"
)

var bytes	=	[16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func (m *Manager) work(q *qt.Queue) []byte{
	l := make([][]byte, 0)
	for t, _ := q.Pop(); t != nil; t, _ = q.Pop(){
		l = append(l, m.next(t).Print())
	}
	switch len(l) {
	case 0:
		return static.ErrWithOutReturn
	case 1:
		return l[0]
	default:
		return m.toPackage(l)
	}
}

func (m *Manager) next(t *qt.Tree) class.Message{
	for n, k := t.GetNode(); n != nil; n, k = t.GetNode(){
		t.SetValue(m.next(n), k)
	}
	l := t.GetValue()
	if len(l) < 2{
		return static.ErrNotComplete
	}
	res := m.tree.Work(l[0].ToString(), l[1].ToString(), l[2:])
	if res.IsWrote{
		m.log.Write(l)
	}
	return res
}

func (m *Manager) toPackage(l [][]byte) []byte{
	s := []byte{'-'}
	length := len(l)

	s = append(s, bytes[length % 16])
	length /= 16

	for length != 0{
		s = append(s, bytes[length % 16])
		length /= 16
	}

	s = append(s, ';')

	for _, v := range l{
		s = append(s, v...)
	}

	return s
}

func (m *Manager) load(q *qt.Queue) []byte{
	l := make([][]byte, 0)
	for t, _ := q.Pop(); t != nil; t, _ = q.Pop(){
		l = append(l, m.nextLoad(t).Print())
	}
	switch len(l) {
	case 0:
		return static.ErrWithOutReturn
	case 1:
		return l[0]
	default:
		return m.toPackage(l)
	}
}

func (m *Manager) nextLoad(t *qt.Tree) class.Message{
	for n, k := t.GetNode(); n != nil; n, k = t.GetNode(){
		t.SetValue(m.next(n), k)
	}
	l := t.GetValue()
	if len(l) < 2{
		return static.ErrNotComplete
	}
	res := m.tree.Work(l[0].ToString(), l[1].ToString(), l[2:])
	return res
}