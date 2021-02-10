package persistence

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateDir() error{
	fmt.Println("create dir [/etc/qlite]")
	return os.Mkdir("/etc/qlite",os.ModePerm)
}

func CreateFile(path string) error{
	path = "/etc/qlite/" + path
	fmt.Println("create:",path)
	if isExist(path){
		return errors.New(path+" is exists")
	}
	f,err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil{
		return err
	}
	defer f.Close()
	return nil
}

func DownloadYaml() error{
	fmt.Println("download qlite.yaml")
	path := "/etc/qlite/qlite.yaml"
	if isExist(path){
		_ = os.Remove(path)
	}
	res, err := http.Get("https://github.com/culion-bear/qlite/releases/download/v2.2.4/qlite.yaml")
	if err != nil {
		return err
	}
	f,err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil{
		return err
	}
	defer f.Close()
	_,err = io.Copy(f,res.Body)
	return err
}

func isExist(path string) bool{
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}