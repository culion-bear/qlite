package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"qlite/hash"
	"qlite/network"
	"qlite/persistence"
	api "qlite/stl"
	"syscall"
)

func toConfig(){
	yamlFile,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err=yaml.Unmarshal(yamlFile,&c)
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig(){
	if len(c.TokenKey) != 16{
		fmt.Println("token key length must be 16")
		os.Exit(1)
	}
	if c.Database <= 0 ||c.Database > 128{
		fmt.Println("database num is error")
		os.Exit(1)
	}
	api.AofHandle = persistence.NewAofHandle(c.AofPath,c.AofInterval,fileWriteError)
	api.LogHandle = persistence.NewLogHandle(c.LogPath,fileWriteError)
	api.StlHandle = persistence.NewStlHandle(c.StlPath,fileWriteError)
	network.Password  = c.Password
	network.TokenKey  = c.TokenKey
	err := hash.HashInit(c.Database)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = api.StlHandle.Read(restoreStl)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = api.AofHandle.Read(restoreDatabase)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}

func fileWriteError(name string,err error){
	if err == syscall.EISDIR {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("[ERROR] FILE:",name,err)
}