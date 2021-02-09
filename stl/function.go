package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func Service(name string) (*stlServiceInfo,error){
	v,flag := servers[name]
	if !flag{
		return nil,ErrNotExist
	}
	return v,nil
}

func GetStlList() []map[string]interface{}{
	stl := make([]map[string]interface{},0)
	for k,v := range servers{
		stl = append(stl, map[string]interface{}{
			"name":k,
			"is_orderly":v.isOrderly,
		})
	}
	return stl
}

func NewService(url string,password string) error{
	conn,err := grpc.Dial(url,grpc.WithInsecure())
	if err != nil{
		return err
	}
	client := NewStlClient(conn)
	serviceInfo,err := client.GetInfo(context.Background(),&Null{})
	if err != nil{
		return err
	}
	if _,flag := servers[serviceInfo.GetName()];flag{
		return ErrServiceExist
	}
	result,err := client.Login(context.Background(),&User{
		Password:             password,
	})
	if err != nil{
		return err
	}
	if result.GetCode() != 200{
		return ErrServiceToken
	}
	apiMap,err := client.GetApiMap(context.Background(),&Null{Token:result.GetToken()})
	if err != nil{
		return err
	}
	clientHandle := &stlServiceInfo{
		url:url,
		handle:client,
		isOrderly:serviceInfo.GetIsOrderly(),
		apiMap:apiMap.GetApi(),
		token:result.GetToken(),
	}
	servers[serviceInfo.GetName()] = clientHandle
	go clientHandle.toHealth()
	return nil
}

func Flush(){
	for _,v := range servers{
		_,_ = v.handle.Flush(context.Background(),&Null{Token:v.token})
	}
}

func (handle *stlServiceInfo) toHealth() {
	for true{
		time.Sleep(time.Second*15)
		_,err := handle.handle.Ping(context.Background(),&Null{})
		if err != nil{
			break
		}
	}
	go handle.restart()
}

func (handle *stlServiceInfo) restart() {
	sum := 0
	for true{
		sum ++
		conn,err := grpc.Dial(handle.url,grpc.WithInsecure())
		if err != nil{
			fmt.Printf("[%s]第%d次重连失败...\n",handle.url,sum)
			time.Sleep(time.Second*5)
		}else{
			handle.handle = NewStlClient(conn)
			_,err = handle.handle.Ping(context.Background(),&Null{})
			if err != nil{
				fmt.Printf("[%s]第%d次重连失败...\n",handle.url,sum)
				time.Sleep(time.Second*5)
			}else{
				fmt.Printf("[%s]第%d次重连成功\n",handle.url,sum)
				go handle.toHealth()
				return
			}
		}
	}
}

func (handle *stlServiceInfo) IsOrderly() bool{
	return handle.isOrderly
}

func (handle *stlServiceInfo) IsExists(key string) bool{
	_,flag := handle.apiMap[key]
	return flag
}

func (handle *stlServiceInfo) IsWriting(key string) bool{
	k,flag := handle.apiMap[key]
	if !flag{
		return false
	}
	return k.GetIsWriting()
}

func (handle *stlServiceInfo) GetApiInfo(key string) (ApiInfo,error){
	k,flag := handle.apiMap[key]
	if !flag{
		return ApiInfo{},ErrApiNotExit
	}
	return *k,nil
}

func (handle *stlServiceInfo) GetApiDescriptionList() ([]*ApiDescription,error){
	result,err := handle.handle.GetApiDescriptionList(context.Background(),&Null{Token:handle.token})
	if err != nil{
		return nil,err
	}
	return result.GetList(),nil
}

func (handle *stlServiceInfo) Submit(msg *PendingMessage) (*Result,error){
	msg.Token = handle.token
	return handle.handle.Submit(context.Background(),msg)
}

func (handle *stlServiceInfo) Compute(msg *Request) (*Response,error){
	msg.Token = handle.token
	return handle.handle.Compute(context.Background(),msg)
}