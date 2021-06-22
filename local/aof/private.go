package aof

import (
	"errors"
	class "qlite/struct"
	"strconv"
	"viv/local/file"
)

var bytes	=	[16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func (m *Manager) validation() error{
	sum := file.GetFileCount(m.path)
	for i := 0; i < sum; i ++{
		fName := m.toFileName(i)
		if !file.FileIsExists(fName){
			return errors.New(fName + " is not found")
		}
	}
	return nil
}

func (m *Manager) toFileName(num int) string{
	return m.path + "/" + strconv.Itoa(num) + ".aof"
}

func (m *Manager) initAof() (*writer, int){
	number := file.GetFileCount(m.path)
	if number == 0{
		return newWriter(m.toFileName(number), m.t), number + 1
	}
	size := file.GetFileSize(m.toFileName(number - 1))
	if size >= m.maxSize{
		return newWriter(m.toFileName(number), m.t), number + 1
	}
	m.size = size
	return newWriter(m.toFileName(number - 1), m.t), number
}

func (m *Manager) isMax() bool{
	return m.size >= m.maxSize
}

func (m *Manager) next(number int, handle *writer) (*writer, int){
	handle.Close()
	handle = newWriter(m.toFileName(number), m.t)
	m.size = 0
	return handle, number + 1
}

func (m *Manager) toPackage(l []class.Message) []byte{
	s := []byte{'='}
	length := len(l)

	s = append(s, bytes[length % 16])
	length /= 16

	for length != 0{
		s = append(s, bytes[length % 16])
		length /= 16
	}

	s = append(s, ';')

	for _, v := range l{
		s = append(s, v.Print()...)
	}

	m.size += len(s)

	return s
}