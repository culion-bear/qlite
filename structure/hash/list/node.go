package list

import "qlite/node"

type List struct {
	last	*List
	next	*List
	value	node.Node
}

func newNode(v node.Node) *List {
	return &List{
		last:  nil,
		next:  nil,
		value: v,
	}
}