package sd

import "viv/structure/qt"

type Model struct {
	msg		*qt.Queue
	ch		chan []byte
}

func NewModel(ch chan []byte, msg *qt.Queue) *Model{
	return &Model{
		ch:  ch,
		msg: msg,
	}
}