package main

import (
	"context"
	"encoding/json"
	"errors"
	flags "flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/signal"
	"qlite/hash"
	"qlite/network"
	"qlite/persistence"
	api "qlite/stl"
	"syscall"
	"time"
)

func init(){
	lTime.Start()
	flags.StringVar(&fileName,"f","/etc/qlite/qlite.yaml","yaml file path")
	flags.BoolVar(&version,"v",false,"show version")
	flags.BoolVar(&help,"h",false,"show help")
	flags.BoolVar(&daemon,"d",false,"to daemon")
	flags.Parse()
}

func toConfig() config{
	yamlFile,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		panic(err)
	}
	c := new(config)
	err=yaml.Unmarshal(yamlFile,c)
	if err!=nil{
		panic(err)
	}
	return *c
}

func initConfig(c config){
	if len(c.TokenKey) != 16{
		panic(errors.New("token key length must be 16"))
	}
	if c.Database <= 0 ||c.Database >= 256{
		panic(errors.New("database num is error"))
	}
	network.AofHandle = persistence.NewAofHandle(c.AofPath,c.AofInterval,fileWriteError)
	network.LogHandle = persistence.NewLogHandle(c.LogPath,fileWriteError)
	network.StlHandle = persistence.NewStlHandle(c.StlPath,fileWriteError)
	network.Password  = c.Password
	network.TokenKey  = c.TokenKey
	err := hash.HashInit(c.Database)
	if err != nil{
		panic(err)
	}
	err = network.StlHandle.Read(restoreStl)
	if err != nil{
		panic(err)
	}
	err = network.AofHandle.Read(restoreDatabase)
	if err != nil{
		panic(err)
	}
}

func fileWriteError(name string,err error){
	if err == syscall.EISDIR {
		panic(err)
	}
	fmt.Println("[ERROR] FILE:",name,err)
}

func restoreStl(url string) bool{
	err := api.NewService(url)
	//if err != nil{
	//	panic(err)
	//}
	if err != nil{
		var r string
		fmt.Printf("[%s ERROR CONNECT](%s) would you want to connect again ? [y/N] :",url,err.Error())
		_, _ = fmt.Scan(&r)
		return r == "y"
	}
	return false
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
		case api.OptionCode_UPDATE:
			//if result.GetNewId() == ""{
			//	fmt.Println("[ERROR]OPTION_CODE_UPDATE:","NEW_ID_IS_EMPTY")
			//	return
			//}
			//err = base.UpdateNodeID(msg.Key,result.GetNewId())
			//if err != nil{
			//	fmt.Println("[ERROR]OPTION_CODE_UPDATE:",err)
			//	return
			//}
		case api.OptionCode_DELETE:
			base.Del([]string{msg.Key})
		default:
			break
		}
	}
}

func kill(handle *iris.Application){
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		os.Interrupt,
		syscall.SIGINT,
		os.Kill,
		syscall.SIGKILL,
		syscall.SIGTERM,
	)
	select {
	case <-ch:
		println("wait...")
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_ = handle.Shutdown(ctx)
		network.AofHandle.Close()
		network.LogHandle.Close()
		network.StlHandle.Close()
		api.Flush()
	}
}

func main(){
	fmt.Println(logo)
	if version||help{
		if version{
			fmt.Println(network.Version)
		}
		if help{
			flags.PrintDefaults()
		}
		return
	}
	c := toConfig()
	initConfig(c)
	handle := network.IrisInit(c.CORS)
	go kill(handle)
	_ = handle.Run(iris.Addr(fmt.Sprintf("%s:%d",c.IP,c.Port)), iris.WithoutInterruptHandler)
}