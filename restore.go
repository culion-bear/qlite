package main

import (
	"encoding/json"
	"fmt"
	"os"
	"qlite/hash"
	"qlite/network"
	"qlite/persistence"
	api "qlite/stl"
)

func restoreStl(buf []byte) interface{}{
	var msg network.StlUrl
	err := json.Unmarshal(buf,&msg)
	if err != nil{
		fmt.Printf("[STL JSON ERROR READ]%s",err.Error())
		os.Exit(1)
	}
	err = api.NewService(msg.Url,msg.Password)
	if err == nil{
		return nil
	}
	ok := false
	for true{
		var r string
		fmt.Printf("[%s ERROR CONNECT](%s) would you want to connect again ? [y/N] :",msg.Url,err.Error())
		_, _ = fmt.Scan(&r)
		if !(r == "y" || r == "Y"){
			return nil
		}
		if ok = err == api.ErrServiceToken;ok{
			fmt.Printf("[%s] scan password again:",msg.Url)
			_,_ = fmt.Scan(&msg.Password)
		}
		err = api.NewService(msg.Url,msg.Password)
		if err == nil{
			break
		}
	}
	if ok{
		return msg
	}
	return nil
}

func restoreDatabase(msg persistence.Data){
	t := lTime.GetTime()
	if msg.Time > 0 && msg.BeginTime + msg.Time < t{
		return
	}
	if msg.Time <= 0{
		t = 0
	}else{
		t = msg.Time + msg.BeginTime - t
	}
	if msg.Type == "hash"{
		var err error
		switch msg.Option {
		case "flush":
			err = hash.Flush(msg.Database)
		case "set":
			err = hash.SetHash(msg.Database,msg.Path,msg.Key,t)
		case "set_x":
			err = hash.SetHashX(msg.Database,msg.Path,msg.Key,t)
		case "delete":
			_,err = hash.Delete(msg.Database,msg.Path,msg.Keys)
		case "pex":
			err = hash.Pex(msg.Database,msg.Path,msg.Key,t)
		case "pex_to":
			err = hash.PexTo(msg.Database,msg.Path,msg.Key,t)
		case "rename":
			err = hash.Rename(msg.Database,msg.Path,msg.Key,msg.NewKey)
		case "rename_x":
			err = hash.RenameX(msg.Database,msg.Path,msg.Key,msg.NewKey)
		default:
			fmt.Println("[DEFAULT]RESTORE:",msg)
		}
		if err != nil{
			fmt.Println("[ERROR]RESTORE:",err)
		}
	}else{
		service,err := api.Service(msg.Type)
		if err != nil{
			fmt.Println("[ERROR]SERVICE:",err)
			return
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
		result,err := service.Submit(&api.PendingMessage{
			Id:                   id,
			Opt:                  msg.Option,
			Msg:                  opt,
			Time:                 msg.BeginTime,
		})
		if err != nil{
			fmt.Println("[ERROR]RESULT:",err)
			return
		}
		switch result.GetCode() {
		case api.OptionCode_CREATE:
			if result.GetNewId() == ""{
				fmt.Println("[ERROR]OPTION_CODE_CREATE:","NEW_ID_IS_EMPTY")
				return
			}
			err = base.Set(hash.NewOtherNode(result.GetNewId(),msg.Key,msg.Type,t))
			if err != nil{
				fmt.Println("[ERROR]OPTION_CODE_CREATE:",err)
				return
			}
		case api.OptionCode_DELETE:
			base.Del([]string{msg.Key})
		default:
			break
		}
	}
}