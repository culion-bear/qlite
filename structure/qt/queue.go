package qt

type Queue struct {
	head	*node
	tail	*node
}

func NewQueue() *Queue{
	return &Queue{}
}

func (m *Queue) Push(v *Tree, p int){
	n := newNode(v, p)
	if m.head == nil{
		m.head, m.tail = n, n
		return
	}
	m.tail.next = n
	m.tail = n
}

func (m *Queue) Pop() (*Tree, int){
	if m.head == nil{
		return nil, 0
	}
	defer m.next()
	return m.head.value, m.head.point
}

func (m *Queue) next(){
	m.head = m.head.next
}