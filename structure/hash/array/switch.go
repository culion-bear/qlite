package array

func (m *Manager) GetCap() int{
	return m.size
}

func (m *Manager) Lessen() {
	if m.size <= 16 || m.getLoadFactor() > 0.1{
		return
	}
	m.rehash(m.size / 2)
}

func (m *Manager) Dilatation() {
	if m.getLoadFactor() < 1.0 || m.size >= 2147483648{
		return
	}
	m.rehash(m.size * 2)
}

func (m *Manager) LessenWithSize(size int) {
	if m.size <= 16 || m.getLoadFactor() > 0.1{
		return
	}
	m.rehash(size)
}

func (m *Manager) GetNextPot(x int) int{
	x = x - 1
	x = x | (x >> 1)
	x = x | (x >> 2)
	x = x | (x >> 4)
	x = x | (x >> 8)
	x = x | (x >>16)
	return x + 1
}