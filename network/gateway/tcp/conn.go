package tcp

import (
	"net"
	"viv/controller/resolver"
	"viv/network/scanner"
	"viv/network/tcp"
)

type Conn struct {
	handle		net.Conn
	chResponse	chan []byte
	chClose		chan bool
	scan		*scanner.Manager
	res			*resolver.Manager
}

func New(conn net.Conn) tcp.Conn{
	c := &Conn{
		handle:     conn,
		chResponse: make(chan []byte),
		chClose:    make(chan bool),
		res:        resolver.New(),
	}
	c.scan = scanner.New(c.option)
	return c
}