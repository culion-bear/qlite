package main

import (
	flags "flag"
	"fmt"
	"os"
	"qlite/network"
	"qlite/persistence"
	"syscall"
)

func init(){
	lTime.Start()
	flags.StringVar(&fileName,"f","/etc/qlite/qlite.yaml","yaml file path")
	flags.BoolVar(&initPath,"i",false,"init config for QLite")
	flags.BoolVar(&version,"v",false,"show version")
	flags.BoolVar(&help,"h",false,"show help")
	flags.BoolVar(&daemon,"d",false,"to daemon")
	flags.StringVar(&signal,"s","","send signal to QLite process: [stop,reload]")
	flags.Parse()
}

func setVersion(v string){
	network.Version = v
}

func initQLite(){
	if !initPath{
		return
	}
	if !isRoot(){
		fmt.Println("please run this initialize process as root")
		os.Exit(1)
	}
	fmt.Println("init config for QLite...")
	err := persistence.CreateDir()
	if err != nil{
		fmt.Println(err)
	}
	err = persistence.DownloadYaml()
	if err != nil{
		fmt.Println(err)
	}
	err = persistence.CreateFile("qlite.aof")
	if err != nil{
		fmt.Println(err)
	}
	err = persistence.CreateFile("qlite.stl")
	if err != nil{
		fmt.Println(err)
	}
	err = persistence.CreateFile("qlite.log")
	if err != nil{
		fmt.Println(err)
	}
	err = persistence.CreateFile("qlite.pid")
	if err != nil{
		fmt.Println(err)
	}
	os.Exit(0)
}

func showInfo(){
	ok := version || help
	if version{
		fmt.Println(network.Version)
	}
	if help{
		flags.PrintDefaults()
	}
	if ok{
		os.Exit(0)
	}
}

func toSignal(pid int){
	switch signal {
	case "stop":
		sendSignal(pid,syscall.SIGTERM)
	case "reload":
		reload(pid)
	default:
		return
	}
	os.Exit(0)
}

func run(){
	initQLite()
	showInfo()
	toConfig()
	pid := getPid()
	toSignal(pid)
	if isRun(pid){
		fmt.Printf("[%d] is running\n",pid)
		os.Exit(1)
	}
	if daemon{
		toDaemon()
	}
	start()
}

func main(){
	setVersion("V 2.2.4")
	fmt.Println(logo)
	run()
}