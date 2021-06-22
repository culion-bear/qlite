package file

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func GetFileCount(path string) int{
	l, _ := ioutil.ReadDir(path)
	return len(l)
}

func GetFileName(path string) ([]string, error){
	l, err := ioutil.ReadDir(path)
	if err != nil{
		return nil, err
	}
	fs, flag := make([]string, len(l)), path[len(path) - 1] == '/'
	for k, v := range l{
		fs[k] = func() string{
			if flag{
				return path + v.Name()
			}
			return path + "/" + v.Name()
		}()
	}
	return fs, nil
}

func FileIsExists(path string) bool{
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GetFileSize(path string) int{
	file, _ := os.Open(path)
	fi,_ := file.Stat()
	return int(fi.Size())
}

func ReadFile(path string) ([]byte, error){
	f, err := os.Open(path)
	if err != nil{
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func WriteFile(path string, msg interface{}) error{
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil{
		return err
	}
	defer f.Close()
	_,err = fmt.Fprintf(f, "%v", msg)
	return err
}

func CreateDir(path string) error{
	return os.Mkdir(path, os.ModePerm)
}

func CreateFile(path string) error{
	if FileIsExists(path){
		return errors.New(path+" is exists")
	}
	f,err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil{
		return err
	}
	return f.Close()
}

func DownloadFile(path, url string) error{
	if FileIsExists(path){
		err := os.Remove(path)
		if err != nil{
			return err
		}
	}
	res, err := http.Get(url)
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