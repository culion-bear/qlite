package scanner

type handle func ([]byte)

type Manager struct {
	size		int
	cache		int
	point		int
	buf			[]byte
	worker		handle
}

func New(f handle) *Manager{
	return &Manager{
		size:     0,
		cache:    0,
		point:    0,
		buf:      make([]byte, 0),
		worker:   f,
	}
}