package tcp

import (
	"bytes"
	sd "viv/controller/scheduler"
	"viv/static"
)

var point = []byte("\r\n")

func (m *Conn) write(buf []byte) error{
	_, err := m.handle.Write(m.scan.Package(buf))
	return err
}

func (m *Conn) option(buf []byte){
	q, err := m.res.Split(buf)
	if err != nil{
		str := []byte{'!'}
		str = append(str, m.scan.Package(err)...)
		m.chResponse <- str
		return
	}
	sd.Scheduler.Option(sd.NewModel(m.chResponse, q))
}

func (m *Conn) getPassword() error{
	buf := make([]byte, 256)
	n, err := m.handle.Read(buf)
	if err != nil{
		return err
	}
	l, err := m.split(buf[: n])
	if err != nil{
		return err
	}
	if !bytes.Equal(l[0], sd.Scheduler.GetPassword()){
		return static.ErrPasswordError
	}
	m.scan.Save(l[1])
	return nil
}

func (m *Conn) split(buf []byte) ([][]byte, error){
	l := bytes.SplitN(buf, point, 2)
	if len(l) != 2{
		return l, static.ErrNotSplitPoint
	}
	return l, nil
}