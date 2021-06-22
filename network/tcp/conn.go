package tcp

import "net"

type newConn func(net.Conn) Conn

type Conn interface {
	Response()
	Listen()
	Close()
}