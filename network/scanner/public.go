package scanner

func (m *Manager) Push(msg []byte) []byte{
	m.buf = append(m.buf, msg...)
	for true{
		res := m.split()
		if !res.isFinish{
			break
		}
		if res.err != nil{
			return res.err
		}
		m.worker(res.msg)
	}
	return nil
}

func (m *Manager) Work(msg []byte) {
	m.worker(msg)
}

func (m *Manager) Package(msg []byte) []byte{
	s := make([]byte, 0)
	length := len(msg)

	s = append(s, bytes[length % 16])
	length /= 16

	for length != 0{
		s = append(s, bytes[length % 16])
		length /= 16
	}

	s = append(s, ';')

	return append(s, msg...)
}

func (m *Manager) Save(msg []byte){
	m.buf = append(m.buf, msg...)
}