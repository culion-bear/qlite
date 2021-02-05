package network

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func init(){
	lTime.Start()
}

func IrisInit(isCORS bool) *iris.Application{
	irisHandle:=iris.New()
	tokenHandle:=initToken()
	var corsHandle router.Party
	if isCORS{
		corsHandle=irisHandle.Party("/cors",cors.New(cors.Options{
			AllowedOrigins:[]string{"*"},
			AllowedMethods:[]string{"POST","GET","OPTIONS","DELETE"},
			MaxAge:3600,
			AllowedHeaders:[]string{"*"},
			AllowCredentials:true,
		})).AllowMethods(iris.MethodOptions)
	}else{
		corsHandle=irisHandle.Party("/")
	}
	corsHandle.Use(setReadLog)
	corsHandle.Post("/login",Login)
	corsHandle.Post("/ping",Ping)
	corsHandle.PartyFunc("/token", func(tokenParty router.Party) {
		tokenParty.Use(tokenHandle.Serve)
		tokenParty.PartyFunc("/service", func(serviceParty router.Party) {
			serviceParty.Get("/info",Info)
			serviceParty.Get("/list/stl",StlList)
			serviceParty.Get("/list/api/{stl}",ApiList)
			serviceParty.Post("/join",Join)
			serviceParty.Get("/database/{num:int}",Database)
			serviceParty.Get("/flush",Flush)
			serviceParty.Get("/number",DatabaseNumber)
		})
		tokenParty.PartyFunc("/hash", func(hashParty router.Party) {
			hashParty.Post("/set",Set)
			hashParty.Post("/set_x",SetX)
			hashParty.Get("/select",Select)
			hashParty.Get("/select_x",SelectX)
			hashParty.Get("/size",Size)
			hashParty.Post("/delete",Delete)
			hashParty.Post("/type",Type)
			hashParty.Post("/exists",Exists)
			hashParty.Post("/pex",Pex)
			hashParty.Post("/pex_to",PexTo)
			hashParty.Post("/time",Time)
			hashParty.Post("/time_to",TimeTo)
			hashParty.Post("/rename",Rename)
			hashParty.Post("/rename_x",RenameX)

			hashParty.Post("/set/{father:path}",Set)
			hashParty.Post("/set_x/{father:path}",SetX)
			hashParty.Get("/select/{father:path}",Select)
			hashParty.Get("/select_x/{father:path}",SelectX)
			hashParty.Get("/size/{father:path}",Size)
			hashParty.Post("/delete/{father:path}",Delete)
			hashParty.Post("/type/{father:path}",Type)
			hashParty.Post("/exists/{father:path}",Exists)
			hashParty.Post("/pex/{father:path}",Pex)
			hashParty.Post("/pex_to/{father:path}",PexTo)
			hashParty.Post("/time/{father:path}",Time)
			hashParty.Post("/time_to/{father:path}",TimeTo)
			hashParty.Post("/rename/{father:path}",Rename)
			hashParty.Post("/rename_x/{father:path}",RenameX)
		})
		tokenParty.PartyFunc("/stl", func(stlParty router.Party) {
			stlParty.Post("/{service}/{api}",StlApi)
			stlParty.Post("/{service}/{api}/{father:path}",StlApi)
		})
		tokenParty.PartyFunc("/alg", func(algParty router.Party) {
			algParty.Post("/{service}/{api}",AlgApi)
		})
	})
	return irisHandle
}