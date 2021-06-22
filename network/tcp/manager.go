package tcp

import (
	"fmt"
	"net"
)

type Manager struct {
	addr		*net.TCPAddr
	listener	*net.TCPListener
}

func New(ip string,port int) (*Manager,error){
	addr,err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ip, port)) //新建tcp连接包
	return &Manager{
		addr: addr,
	},err
}

func (m *Manager) Close() error{
	return m.listener.Close() //关闭tcp连接
}

func (m *Manager) Listen(f newConn) error{
	var err error
	m.listener,err = net.ListenTCP("tcp", m.addr) //监听地址
	if err != nil{
		return err
	}
	m.run(f) //开始轮询获取接入信息
	return nil
}

func (m *Manager) run(f newConn){
	for true{
		conn,err := m.listener.Accept() //获取接入的连接
		if err != nil{
			continue
		}
		go m.chat(f(conn))
	}
}

func (m *Manager) chat(handle Conn){
	defer func() {
		handle.Close()
	}()
	go handle.Response()
	handle.Listen()
}