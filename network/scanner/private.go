package scanner

import (
	"viv/static"
)

type result struct {
	msg			[]byte
	isFinish	bool
	err			[]byte
}

var scale	=	[9]int{0x1, 0x10, 0x100, 0x1000, 0x10000, 0x100000, 0x1000000, 0x10000000, 0x100000000}

var bytes	=	[16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

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

func (m *Manager) split() result{
	if m.size == 0{
		isFinish, err := m.toNumber()
		if !isFinish || err != nil{
			return result{isFinish: isFinish, err: err}
		}
	}
	return m.toMsg()
}

func (m *Manager) toNumber() (bool, []byte){
	for k, v := range m.buf{
		if v == ';'{
			m.size = m.cache
			m.point = 0
			m.cache = 0
			m.buf = m.buf[k + 1: ]
			return true, nil
		}
		if data[v] == -1 || m.point == 9{
			return true, static.ErrNumber
		}
		m.cache += data[v] * scale[m.point]
		m.point ++
	}
	m.buf = []byte{}
	return false, nil
}

func (m *Manager) toMsg() result{
	size := len(m.buf)
	if m.size > size{
		return result{isFinish: false}
	}
	defer func() {
		m.buf = m.buf[m.size: ]
		m.size = 0
	}()
	return result{
		msg:      m.buf[: m.size],
		isFinish: true,
		err:      nil,
	}
}

//func (m *Manager) copy(msg []byte) []byte{
//	buf := make([]byte, len(msg))
//	copy(buf, msg)
//	return buf
//}