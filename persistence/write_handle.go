package persistence

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type WriteManager struct {
	msgChan   chan string
	errChan   chan error
	flushChan chan bool
	closeChan chan bool
	fileName  string
	interval  int
}

func NewWriter(fileName string,intervalTime int) *WriteManager{
	return &WriteManager{
		msgChan:   make(chan string,1024),
		errChan:   make(chan error),
		flushChan: make(chan bool),
		closeChan: make(chan bool),
		fileName:  fileName,
		interval:  intervalTime,
	}
}

func (handle *WriteManager) Start(p func(string,error)){
	go handle.Run()
	go handle.Error(p)
	if handle.interval > 0 {
		go handle.Flush()
	}
}

func (handle *WriteManager) Restart(){
	go handle.Run()
	if handle.interval > 0 {
		go handle.Flush()
	}
}

func (handle *WriteManager) Write(str string) {
	handle.msgChan <- str
}

func (handle *WriteManager) Error(p func(string,error)){
	for true{
		err := <- handle.errChan
		go p(handle.fileName,err)
		time.Sleep(time.Second*5)
		handle.Restart()
	}
}

func (handle *WriteManager) Close(){
	handle.closeChan <- true
}

func (handle *WriteManager) Flush(){
	timer := time.NewTimer(time.Duration(handle.interval)*time.Second)
	for true{
		select {
		case <-handle.closeChan:
			handle.flushChan<-false
			return
		case <-timer.C:
			handle.flushChan<-true
			timer.Reset(time.Duration(handle.interval)*time.Second)
		}
	}
}

func (handle *WriteManager) Run(){
	file,err := os.OpenFile(handle.fileName,os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	if err != nil{
		handle.errChan<-err
		return
	}
	if handle.interval > 0{
		w := bufio.NewWriter(file)
		for true{
			select {
			case msg := <- handle.msgChan:
				_,err = w.WriteString(msg)
				if err != nil{
					handle.errChan<-err
				}
			case flag := <- handle.flushChan:
				if flag{
					err = w.Flush()
					if err != nil{
						handle.errChan<-err
					}
				}else{
					_ = w.Flush()
					_ = file.Close()
					return
				}
			}
		}
	}else{
		for true{
			select {
			case msg := <- handle.msgChan:
				_,err = fmt.Fprintf(file,msg)
				if err != nil{
					handle.errChan<-err
				}
			case <- handle.closeChan:
				_ = file.Close()
				return
			}
		}
	}
}