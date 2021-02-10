package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func ReadPid(path string) (int,error){
	buf,err := ioutil.ReadFile(path)
	if err != nil || len(buf) == 0{
		return -1, err
	}
	return strconv.Atoi(string(buf))
}

func WritePid(path string,pid int) error{
	f,err:=os.OpenFile(path,os.O_WRONLY,os.ModePerm)
	if err!=nil{
		return err
	}
	defer f.Close()
	_,err = f.Write([]byte(strconv.Itoa(pid)))
	return err
}

func GetCmd(pid int) (string,error){
	buf,err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline",pid))
	return string(buf),err
}