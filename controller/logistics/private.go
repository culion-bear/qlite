package logistics

func (m *Manager) option(buf []byte){
	n, err := m.res.Split(buf)
	if err != nil{
		return
	}
	m.load(n)
}