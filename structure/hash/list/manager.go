package list

type Manager struct {
	size	int
	head	*List
	tail	*List
}

func New() *Manager{
	return &Manager{
		size: 0,
		head: nil,
		tail: nil,
	}
}