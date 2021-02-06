package network

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"math/rand"
	"qlite/hash"
	"qlite/persistence"
	"qlite/stl"
	"time"
)

func init(){
	rand.Seed(time.Now().Unix())
}

func setReadLog(ctx iris.Context){
	LogHandle.Write(ctx.Path(),persistence.INFO,persistence.READ)
	ctx.Next()
}

func sendMessage(ctx iris.Context,value interface{},level string){
	LogHandle.Write(ctx.Path(),level,persistence.SEND)
	_,_ = ctx.JSON(value)
}

func sendNormalMessage(ctx iris.Context,value interface{}){
	sendMessage(ctx,value,persistence.INFO)
}

func sendJson(ctx iris.Context,value []byte){
	LogHandle.Write(ctx.Path(),persistence.INFO,persistence.SEND)
	_,_ = ctx.Write(value)
}

func sendError(ctx iris.Context,code int,msg string){
	sendMessage(ctx,Error{
		Code: code,
		Msg:  msg,
	},persistence.ERROR)
}

func Login(ctx iris.Context){
	var msg User
	if err := ctx.ReadJSON(&msg,irisJson);err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.Password != Password{
		sendError(ctx,ErrLogin,"login error")
		return
	}
	token,err := tokenMessage(0)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
		"token":token,
	})
}

func Ping(ctx iris.Context){
	var msg User
	if err := ctx.ReadJSON(&msg,irisJson);err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.Password != Password{
		sendError(ctx,ErrPing,"")
		return
	}
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
	})
}

func Info(ctx iris.Context){
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
		"version":Version,
	})
}

func StlList(ctx iris.Context){
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
		"list":api.GetStlList(),
	})
}

func ApiList(ctx iris.Context){
	s,err := api.Service(ctx.Params().GetString("stl"))
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	l,err := s.GetApiDescriptionList()
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
		"list":l,
	})
}

func Join(ctx iris.Context){
	var msg StlUrl
	if err := ctx.ReadJSON(&msg,irisJson);err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	err := api.NewService(msg.Url)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
	})
	go StlHandle.Write(msg.Url)
}

func Database(ctx iris.Context){
	num,err := ctx.Params().GetInt("num")
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if err = hash.IsOverFlow(num);err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	token,err := tokenMessage(int32(num))
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
		"token":token,
	})
}

func Flush(ctx iris.Context){
	num := getTokenNumber(ctx)
	if err := hash.Flush(num);err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
	})
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "flush",
		Database: num,
		Path:     "",
	})
}

func DatabaseNumber(ctx iris.Context){
	sendNormalMessage(ctx,map[string]interface{}{
		"code":Success,
		"num":hash.GetLength(),
	})
}

func Set(ctx iris.Context){
	var msg TimeModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.Key == ""{
		sendError(ctx,ErrJson,ErrKeyEmpty.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	//if msg.Type == "hash"{
	//	err = hash.SetHash(num,path,msg.Key,msg.Time)
	//}else{
	//	_,err = api.Service(msg.Key)
	//	if err != nil{
	//		sendError(ctx,ErrService,err.Error())
	//		return
	//	}
	//	err = hash.SetOther(num,path,msg.Key,msg.Type,"",msg.Time)
	//}
	err = hash.SetHash(num,path,msg.Key,msg.Time)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "set",
		Database: num,
		Path:     path,
		Key:      msg.Key,
		BeginTime:lTime.GetTime(),
		Time:     msg.Time,
	})
}

func SetX(ctx iris.Context){
	var msg TimeModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.Key == ""{
		sendError(ctx,ErrJson,ErrKeyEmpty.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	//if msg.Type == "hash"{
	//	err = hash.SetHashX(num,path,msg.Key,msg.Time)
	//}else{
	//	_,err = api.Service(msg.Key)
	//	if err != nil{
	//		sendError(ctx,ErrService,err.Error())
	//		return
	//	}
	//	err = hash.SetOtherX(num,path,msg.Key,msg.Type,"",msg.Time)
	//}
	err = hash.SetHashX(num,path,msg.Key,msg.Time)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "set_x",
		Database: num,
		Path:     path,
		Key:      msg.Key,
		BeginTime:lTime.GetTime(),
		Time:     msg.Time,
	})
}

func Select(ctx iris.Context){
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	result,err := hash.Select(num,path)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"list":result,
	},persistence.INFO)
}

func SelectX(ctx iris.Context){
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	result,err := hash.SelectX(num,path)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"list":result,
	},persistence.INFO)
}

func Size(ctx iris.Context){
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	result,err := hash.Size(num,path)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"num":result,
	},persistence.INFO)
}

func Delete(ctx iris.Context){
	var msg KeyList
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	sum, err := hash.Delete(num,path,msg.Keys)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"num":sum,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "delete",
		Database: num,
		Path:     path,
		Keys:     msg.Keys,
	})
}

func Type(ctx iris.Context){
	var msg KeyModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	t,err := hash.Type(num,path,msg.Key)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"type":t,
	},persistence.INFO)
}

func Exists(ctx iris.Context){
	var msg KeyModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	flag,err := hash.Exists(num,path,msg.Key)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"exists":flag,
	},persistence.INFO)
}

func Pex(ctx iris.Context){
	var msg TimeModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	err = hash.Pex(num,path,msg.Key,msg.Time)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "pex",
		Database: num,
		Path:     path,
		Key:      msg.Key,
		BeginTime:lTime.GetTime(),
		Time:     msg.Time,
	})
}

func PexTo(ctx iris.Context){
	var msg TimeModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	err = hash.PexTo(num,path,msg.Key,msg.Time)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "pex_to",
		Database: num,
		Path:     path,
		Key:      msg.Key,
		BeginTime:lTime.GetTime(),
		Time:     msg.Time,
	})
}

func Time(ctx iris.Context){
	var msg KeyModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	t,err := hash.Time(num,path,msg.Key)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"time":t,
	},persistence.INFO)
}

func TimeTo(ctx iris.Context){
	var msg KeyModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	t,err := hash.TimeTo(num,path,msg.Key)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
		"time":t,
	},persistence.INFO)
}

func Rename(ctx iris.Context){
	var msg RenameModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.NewKey == ""{
		sendError(ctx,ErrJson,ErrKeyEmpty.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	err = hash.Rename(num,path,msg.Key,msg.NewKey)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "rename",
		Database: num,
		Path:     path,
		Key:      msg.Key,
		NewKey:   msg.NewKey,
	})
}

func RenameX(ctx iris.Context){
	var msg RenameModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.NewKey == ""{
		sendError(ctx,ErrJson,ErrKeyEmpty.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	err = hash.RenameX(num,path,msg.Key,msg.NewKey)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendMessage(ctx,map[string]interface{}{
		"code":Success,
	},persistence.INFO)
	go AofHandle.Write(persistence.Data{
		Type:     "hash",
		Option:   "rename_x",
		Database: num,
		Path:     path,
		Key:      msg.Key,
		NewKey:   msg.NewKey,
	})
}

func StlApi(ctx iris.Context){
	var msg StlApiModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	if msg.Key == ""{
		sendError(ctx,ErrJson,ErrKeyEmpty.Error())
		return
	}
	var id string
	serviceName,apiName := ctx.Params().Get("service"),ctx.Params().Get("api")
	service,err := api.Service(serviceName)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	info,err := service.GetApiInfo(apiName)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	if info.GetType() != api.STL{
		sendError(ctx,ErrService,api.ErrApiType.Error())
		return
	}
	num,path := getTokenNumber(ctx),ctx.Params().Get("father")
	base,err := hash.ToNode(num,path)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	id, err = base.GetNodeID(msg.Key)
	if err != nil{
		if !info.GetIsWriting(){
			sendError(ctx,ErrService,err.Error())
			return
		}
	}else{
		if t,_ := base.GetNodeType(msg.Key);t != serviceName{
			sendError(ctx,ErrService,api.ErrServiceType.Error())
			return
		}
	}
	opt,_ := json.Marshal(&msg.Opt)
	result,err := service.Submit(&api.PendingMessage{
		Id:                   id,
		Opt:                  apiName,
		Msg:                  opt,
		Time:                 lTime.GetTime(),
	})
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	if result.GetCode() != api.OptionCode_NO_ACTION{
		switch result.GetCode() {
		case api.OptionCode_CREATE:
			if result.GetNewId() == "" {
				sendError(ctx, ErrService, api.ErrIDEmpty.Error())
				return
			}
			err = base.Set(hash.NewOtherNode(result.GetNewId(), msg.Key, serviceName, msg.Time))
			if err != nil {
				sendError(ctx, ErrService, err.Error())
				return
			}
		case api.OptionCode_UPDATE:
			//if result.GetNewId() == ""{
			//	sendError(ctx,ErrService,api.ErrIDEmpty.Error())
			//	return
			//}
			//err = base.UpdateNodeID(msg.Key,result.GetNewId())
			//if err != nil{
			//	sendError(ctx,ErrService,err.Error())
			//	return
			//}
		case api.OptionCode_DELETE:
			base.Del([]string{msg.Key})
		}
		go AofHandle.Write(persistence.Data{
			Type:     serviceName,
			Option:   apiName,
			Database: num,
			Path:     path,
			Detail:   msg.Opt,
			Key:      msg.Key,
			BeginTime:lTime.GetTime(),
			Time:     msg.Time,
		})
	}
	sendJson(ctx,result.GetMsg())
}

func AlgApi(ctx iris.Context){
	var msg AlgApiModel
	err := ctx.ReadJSON(&msg,irisJson)
	if err != nil{
		sendError(ctx,ErrJson,err.Error())
		return
	}
	num := getTokenNumber(ctx)
	serviceName,apiName := ctx.Params().Get("service"),ctx.Params().Get("api")
	service,err := api.Service(serviceName)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	info,err := service.GetApiInfo(apiName)
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	if info.GetType() != api.ALG{
		sendError(ctx,ErrService,api.ErrApiType.Error())
		return
	}
	idList := make([]string,len(msg.Keys))
	for k,v := range msg.Keys{
		id,err := hash.GetID(num,v.Path,v.Key)
		if err != nil{
			sendError(ctx,ErrService,err.Error())
			return
		}
		idList[k] = id
	}
	opt,_ := json.Marshal(msg.Opt)
	result,err := service.Compute(&api.Request{
		Id:                   idList,
		Opt:                  apiName,
		Msg:                  opt,
	})
	if err != nil{
		sendError(ctx,ErrService,err.Error())
		return
	}
	sendJson(ctx,result.GetMsg())
}