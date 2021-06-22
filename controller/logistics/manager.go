package logistics

import (
	"errors"
	class "qlite/struct"
	"viv/controller/resolver"
	"viv/local/aof"
	"viv/structure/qt"
)

type Manager struct {
	writer		*aof.Manager
	ch			chan []class.Message
	load		func(*qt.Queue)
	res			*resolver.Manager
}

func New(path string, t int, f func(*qt.Queue)) (*Manager, error){
	if t < 1{
		return nil, errors.New("interval time must bigger than 1 second")
	}
	m := &Manager{
		writer: nil,
		ch:     make(chan []class.Message, 100000),
		load:   f,
		res:    resolver.New(),
	}
	var err error
	m.writer, err = aof.New(path, t)
	return m, err
}