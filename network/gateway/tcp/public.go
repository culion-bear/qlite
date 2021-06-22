package tcp

import "fmt"

func (m *Conn) Close() {
	m.chClose <- true
	_ = m.handle.Close()
}

func (m *Conn) Response(){
	var msg []byte
	var err error
	for true{
		select {
		case msg = <- m.chResponse:
			err = m.write(msg)
			if err != nil{
				_ = m.handle.Close()
				return
			}
		case <- m.chClose:
			return
		}
	}
}

func (m *Conn) Listen(){
	err := m.getPassword()
	if err != nil{
		fmt.Println(err)
		m.Close()
		return
	}
	errBuf := m.scan.Push([]byte{})
	if errBuf != nil{
		m.Close()
		return
	}
	buf := make([]byte,1024)
	for true{
		n, err := m.handle.Read(buf)
		if err != nil{
			break
		}
		errBuf := m.scan.Push(buf[: n])
		if errBuf != nil{
			break
		}
	}
	m.Close()
}