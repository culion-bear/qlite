package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"os"
	"os/exec"
	slot "os/signal"
	"qlite/network"
	"qlite/persistence"
	api "qlite/stl"
	"strings"
	"syscall"
	"time"
)

func kill(handle *iris.Application){
	ch := make(chan os.Signal, 1)
	slot.Notify(
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
		api.AofHandle.Close()
		api.LogHandle.Close()
		api.StlHandle.Close()
		api.Flush()
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_ = handle.Shutdown(ctx)
	}
}

func start(){
	initConfig()
	setPid()
	handle := network.IrisInit(c.CORS)
	go kill(handle)
	_ = handle.Run(iris.Addr(fmt.Sprintf("%s:%d",c.IP,c.Port)), iris.WithoutInterruptHandler)
}

func toDaemon(){
	//判断当其是否是子进程，当父进程结束之后，子进程会被系统1号进程接管
	if os.Getppid() != 1{
		cmd := exec.Command(os.Args[0],os.Args[1:]...)
		_ = cmd.Start()
		os.Exit(0)
	}
}

func getPid() int{
	pid,err := persistence.ReadPid(c.PidPath)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	return pid
}

func setPid(){
	err := persistence.WritePid(c.PidPath,os.Getpid())
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}

func isRun(pid int) bool{
	if pid == -1{
		return false
	}
	c,err := getCmd(pid)
	return err == nil && len(c) > 0
}

func sendSignal(pid int,msg syscall.Signal){
	p,err := os.FindProcess(pid)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	_ = p.Signal(msg)
}

func getCmd(pid int) (string,error){
	return persistence.GetCmd(pid)
}

func reload(pid int){
	c,err := getCmd(pid)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	args := strings.Split(c,string([]byte{'\x00'}))
	ok := true
	for _,v := range args{
		if v == "-d"{
			ok = false
			break
		}
	}
	if ok{
		args = append(args, "-d")
	}
	sendSignal(pid,syscall.SIGTERM)
	cmd := exec.Command(args[0],args[1:]...)
	_ = cmd.Start()
}

func isRoot() bool{
	return os.Geteuid() == 0
}