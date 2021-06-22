package aof

import (
	"io/ioutil"
	"os"
	class "qlite/struct"
	"viv/local/file"
)

type Manager struct {
	ch		chan []class.Message
	chClose	chan bool
	path	string
	t       int
	maxSize int
	size    int
}

func New(path string, t int) (*Manager, error){
	m := &Manager{
		ch:      make(chan []class.Message, 100000),
		chClose: make(chan bool),
		path:    path,
		maxSize: 4194304,
		t:       t,
		size:    0,
	}
	return m, m.validation()
}

func (m *Manager) Read(fp func([]byte)) error{
	number := file.GetFileCount(m.path)
	for i := 0; i < number; i ++{
		f, err := os.Open(m.toFileName(i))
		if err != nil{
			return err
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil{
			return err
		}
		fp(buf)
	}
	return nil
}

func (m *Manager) Close(){
	m.chClose <- true
}

func (m *Manager) Run(){
	var msg []class.Message
	handle , number := m.initAof()
	handle.Start()
	for true{
		select {
		case msg = <- m.ch:
			if m.isMax(){
				handle, number = m.next(number, handle)
			}
			handle.Write(m.toPackage(msg))
		case <- m.chClose:
			handle.Close()
			return
		}
	}
}

func (m *Manager) Write(msg []class.Message){
	m.ch <- msg
}