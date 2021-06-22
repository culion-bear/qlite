package resolver

import (
	class "qlite/struct"
	"viv/static"
	"viv/structure/qt"
)

var scale	=	[9]int{0x1, 0x10, 0x100, 0x1000, 0x10000, 0x100000, 0x1000000, 0x10000000, 0x100000000}

var data	=	func () [256]int{
	var l [256]int
	for i := 0;i < 256; i ++{
		if i <= '9' && '0' <= i{
			l[i] = i - '0'
		}else if i <= 'f' && 'a' <= i{
			l[i] = i - 'a' + 10
		}else if i <= 'F' && 'A' <= i{
			l[i] = i - 'A' + 10
		}else{
			l[i] = -1
		}
	}
	return l
}()

func (m *Manager) split(buf []byte) (*qt.Queue, static.Error){
	length, q := len(buf), qt.NewQueue()

	var t *qt.Tree
	var err static.Error

	for i := 0; i < length; {
		if buf[i] != static.ResTypeChannel{
			return q, static.ErrTypeNotExist
		}
		t, i, err = m.splitChannel(buf, len(buf), i + 1)
		if err != nil{
			return q, err
		}
		q.Push(t, 0)
	}
	return q, nil
}

func (m *Manager) splitChannel(buf []byte, length, key int) (*qt.Tree, int, static.Error){
	num, key, err := m.toNumber(buf, length, key)
	if err != nil{
		return nil, 0, err
	}

	t := qt.NewTree()
	var str 	[]byte
	var number 	int
	var l 		[][]byte
	var n 		*qt.Tree

	for i := 0; i < num; i ++{
		if key >= length{
			return t, 0, static.ErrPackage
		}
		switch buf[key] {
		case static.ResTypeString, static.ResTypeError, static.ResTypeNode:
			str, key, err = m.splitString(buf, length, key + 1)
			if err != nil{
				return t, 0, err
			}
			t.PushValue(class.Message{
				Type:    static.ResTypeString,
				String:  str,
			})
		case static.ResTypeNumber:
			number, key, err = m.toNumber(buf, length, key + 1)
			if err != nil{
				return t, 0, err
			}
			t.PushValue(class.Message{
				Type:    static.ResTypeNumber,
				Number:  number,
			})
		case static.ResTypeFunc:
			n, key, err = m.splitChannel(buf, length, key + 1)
			if err != nil{
				return t, 0, err
			}
			t.PushNode(n)
		case static.ResTypeList:
			l, key, err = m.splitList(buf, length, key + 1)
			if err != nil{
				return t, 0, err
			}
			t.PushValue(class.Message{
				Type:    static.ResTypeList,
				List:    l,
			})
		default:
			return t, 0, static.ErrTypeNotExist
		}
	}
	return t, key, nil
}

func (m *Manager) splitString(buf []byte, length, key int) ([]byte, int, static.Error){
	num, key, err := m.toNumber(buf, length, key)
	if err != nil{
		return nil, 0, err
	}
	if key + num > length{
		return nil, 0, static.ErrPackage
	}
	return buf[key : key + num], key + num, nil
}

func (m *Manager) splitList(buf []byte, length, key int) ([][]byte, int, static.Error){
	num, key, err := m.toNumber(buf, length, key)
	if err != nil{
		return nil, 0, err
	}
	l := make([][]byte, num)
	var str []byte
	for i := 0; i < num; i ++{
		str, key, err = m.splitString(buf, length, key)
		if err != nil{
			return nil, 0, err
		}
		l[i] = str
	}
	return l, key, nil
}

func (m *Manager) toNumber(buf []byte, length, key int) (int, int, static.Error){
	num := 0
	for i := key; i < length; i ++{
		if buf[i] == ';'{
			return num, i + 1, nil
		}
		if data[buf[i]] == -1 || i - key == 9{
			return 0, 0, static.ErrNumber
		}
		num += data[buf[i]] * scale[i - key]
	}
	return 0, 0, static.ErrNotResPoint
}