package persistence

import (
	"bufio"
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

func (handle *StlManager) Write(url string){
	handle.writer.Write(url+"\n")
}

func (handle *StlManager) Read(p func(string)) error{
	f,err:=os.Open(handle.path)
	if err!=nil{
		return err
	}
	defer f.Close()
	rf := bufio.NewReader(f)
	for true{
		buf,_,err := rf.ReadLine()
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}
		p(string(buf))
	}
	return nil
}