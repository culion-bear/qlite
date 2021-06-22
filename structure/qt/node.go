package qt

type node struct {
	value	*Tree
	point	int
	next	*node
}

func newNode(v *Tree, p int) *node{
	return &node{
		value: v,
		point: p,
		next:  nil,
	}
}