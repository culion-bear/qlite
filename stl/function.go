package api

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"qlite/hash"
	"qlite/persistence"
	"time"
)

func init(){
	lTime.Start()
}

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
		url:		url,
		name:		serviceInfo.GetName(),
		handle:		client,
		isOrderly:	serviceInfo.GetIsOrderly(),
		apiMap:		apiMap.GetApi(),
		token:		result.GetToken(),
		flag:		false,
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
		if err == nil{
			handle.handle = NewStlClient(conn)
			_,err = handle.handle.Ping(context.Background(),&Null{})
			if err == nil{
				fmt.Printf("(%s)[%s]reconnect success - %d\n",handle.name,handle.url,sum)
				handle.Restore()
				go handle.toHealth()
				return
			}
		}
		fmt.Printf("(%s)[%s]reconnect failed - %d\n",handle.name,handle.url,sum)
		time.Sleep(time.Second*5)
	}
}

func (handle *stlServiceInfo) Restore(){
	handle.lock()
	defer handle.unlock()
	flag,err := handle.handle.Exists(context.Background(),&User{
		Password:             handle.token,
	})
	if err != nil || flag.GetOk(){
		return
	}
	AofHandle.Flush()
	l,err := AofHandle.Restore(handle.name)
	if err != nil{
		fmt.Printf("(%s)[%s]database file was read failed：%s\n",handle.name,handle.url,err.Error())
		return
	}
	for _,v := range l{
		handle.restoreData(v)
	}
}

func (handle *stlServiceInfo) restoreData(msg persistence.Data){
	t := lTime.GetTime()
	if msg.Time > 0 && msg.BeginTime + msg.Time < t{
		return
	}
	if msg.Time <= 0{
		t = 0
	}else{
		t = msg.Time + msg.BeginTime - t
	}
	base,err := hash.ToNode(msg.Database,msg.Path)
	if err != nil{
		fmt.Println("[ERROR]DATABASE:",err)
		return
	}
	id,err := base.GetNodeID(msg.Key)
	if err != nil{
		fmt.Println("[WAINING]GET_NODE_ID:",err)
	}
	opt,_ := json.Marshal(&msg.Detail)
	result,err := handle.handle.Submit(context.Background(),&PendingMessage{
		Id:                   id,
		Opt:                  msg.Option,
		Msg:                  opt,
		Time:                 msg.BeginTime,
		Token:				  handle.token,
	})
	if err != nil{
		fmt.Println("[ERROR]RESULT:",err)
		return
	}
	switch result.GetCode() {
	case OptionCode_CREATE:
		if result.GetNewId() == ""{
			fmt.Println("[ERROR]OPTION_CODE_CREATE:","NEW_ID_IS_EMPTY")
			return
		}
		base.SetX(hash.NewOtherNode(result.GetNewId(),msg.Key,msg.Type,t))
	case OptionCode_DELETE:
		base.Del([]string{msg.Key})
	default:
		break
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
	if handle.isLock(){
		return &Result{},ErrRestore
	}
	msg.Token = handle.token
	return handle.handle.Submit(context.Background(),msg)
}

func (handle *stlServiceInfo) Compute(msg *Request) (*Response,error){
	if handle.isLock(){
		return &Response{},ErrRestore
	}
	msg.Token = handle.token
	return handle.handle.Compute(context.Background(),msg)
}

func (handle *stlServiceInfo) lock(){
	handle.mu.Lock()
	defer handle.mu.Unlock()
	handle.flag = true
}

func (handle *stlServiceInfo) unlock(){
	handle.mu.Lock()
	defer handle.mu.Unlock()
	handle.flag = false
}

func (handle *stlServiceInfo) isLock() bool{
	handle.mu.RLock()
	defer handle.mu.RUnlock()
	return handle.flag
}