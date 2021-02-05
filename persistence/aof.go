package persistence

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type Data struct {
	Type     	string			`json:"type"`
	Option   	string			`json:"option"`
	Database 	int				`json:"database"`
	Path     	string 			`json:"path"`
	Detail   	interface{}		`json:"detail,omitempty"`
	Key      	string			`json:"key,omitempty"`
	Keys		[]string		`json:"keys,omitempty"`
	NewKey		string			`json:"new_key,omitempty"`
	Time     	int64			`json:"time,omitempty"`
	BeginTime	int64			`json:"begin_time,omitempty"`
}

type AofManager struct {
	path string
	writer *WriteManager
}

func NewAofHandle(fileName string,intervalTime int,p func(string,error)) *AofManager{
	v := &AofManager{
		path:  fileName,
		writer:NewWriter(fileName,intervalTime),
	}
	v.writer.Start(p)
	return v
}

func (handle *AofManager) Close(){
	handle.writer.Close()
}

func (handle *AofManager) Write(msg Data){
	buf,_ := json.Marshal(&msg)
	handle.writer.Write(string(buf)+"\n")
}

func (handle *AofManager) Read(p func(Data)) error{
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
		var v Data
		err = json.Unmarshal(buf,&v)
		if err != nil{
			return err
		}
		p(v)
	}
	return nil
}