package logistics

import class "qlite/struct"

func (m *Manager) Load() error{
	return m.writer.Read(m.option)
}

func (m *Manager) Close(){
	m.ch <- nil
}

func (m *Manager) Run(){
	var msg []class.Message
	go m.writer.Run()
	for true{
		msg = <- m.ch
		if msg == nil{
			m.writer.Close()
			return
		}
		m.writer.Write(msg)
	}
}

func (m *Manager) Write(msg []class.Message){
	m.ch <- msg
}