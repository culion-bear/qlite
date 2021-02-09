package persistence

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type StlManager struct {
	path string
	writer *WriteManager
}

func NewStlHandle(fileName string,p func(string,error)) *StlManager{
	v := &StlManager{
		path:  fileName,
		writer:NewWriter(fileName,0),
	}
	v.writer.Start(p)
	return v
}

func (handle *StlManager) Close(){
	handle.writer.Close()
}

func (handle *StlManager) Write(v interface{}){
	buf,_ := json.Marshal(v)
	handle.writer.Write(string(buf)+"\n")
}

func (handle *StlManager) Read(p func([]byte) interface{}) error{
	f,err:=os.OpenFile(handle.path,os.O_RDWR,os.ModePerm)
	if err!=nil{
		return err
	}
	defer f.Close()
	rf := bufio.NewReader(f)
	ok := false
	msg := make([][]byte,0)
	for true{
		buf,_,err := rf.ReadLine()
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}
		v := p(buf)
		if ok = v != nil;ok{
			buf,_ = json.Marshal(v)
		}
		msg = append(msg, buf)
	}
	if ok{
		for _,v := range msg{
			_,_ = f.WriteString(string(v)+"\n")
		}
	}
	return nil
}