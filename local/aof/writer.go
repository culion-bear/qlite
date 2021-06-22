package aof

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type writer struct {
	msgChan   chan []byte
	errChan   chan error
	closeChan chan bool
	fileName  string
	interval  int
}

func newWriter(fileName string,intervalTime int) *writer{
	return &writer{
		msgChan:   make(chan []byte, 100000),
		errChan:   make(chan error),
		closeChan: make(chan bool),
		fileName:  fileName,
		interval:  intervalTime,
	}
}

func (handle *writer) Start(){
	go handle.run()
	go handle.error()
	go handle.flush()
}

func (handle *writer) Write(str []byte) {
	handle.msgChan <- str
}

func (handle *writer) error(){
	for true{
		err := <- handle.errChan
		if err == nil{
			return
		}
		fmt.Println(handle.fileName,"[aof error] :",err)
	}
}

func (handle *writer) Close(){
	handle.closeChan <- true
	handle.errChan <- nil
}

func (handle *writer) flush(){
	timer := time.NewTimer(time.Duration(handle.interval) * time.Second)
	msg := make([]byte, 0 , 0)
	for true{
		select {
		case <- handle.closeChan:
			handle.msgChan <- nil
			return
		case <- timer.C:
			handle.msgChan <- msg
			timer.Reset(time.Duration(handle.interval) * time.Second)
		}
	}
}

func (handle *writer) toFlush(){
	handle.msgChan <- []byte{}
}

func (handle *writer) run(){
	file,err := os.OpenFile(handle.fileName,os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	if err != nil{
		handle.errChan<-err
		return
	}
	handle.write(file)
}

func (handle *writer) write(file *os.File){
	var err error
	w := bufio.NewWriter(file)
	for true{
		msg := <- handle.msgChan
		if msg == nil{
			_ = w.Flush()
			_ = file.Close()
			return
		}else if len(msg) == 0{
			err = w.Flush()
			if err != nil{
				handle.errChan <- err
			}
		}else{
			_,err = w.Write(msg)
			if err != nil{
				handle.errChan <- err
			}
		}
	}
}