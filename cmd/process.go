package cmd

import (
	"fmt"
	"os"
	"os/exec"
	slot "os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	sd "viv/controller/scheduler"
	"viv/local/file"
	"viv/network/tcp"
)

func kill(m *tcp.Manager) {
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
		err := m.Close()
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
		sd.Scheduler.Close()
		time.Sleep(time.Second)
		os.Exit(0)
	}
}

func toDaemon() error{
	args := make([]string, 0)
	for i := 1; i < len(os.Args); i ++{
		if os.Args[i] != "-d"{
			args = append(args, os.Args[i])
		}
	}
	cmd := exec.Command(os.Args[0],args...)
	return cmd.Start()
}

func isChildProcess() bool{
	//判断当其是否是子进程，当父进程结束之后，子进程会被系统1号进程接管
	return os.Getppid() == 1
}

func getPid(path string) (int, error){
	buf, err := file.ReadFile(path)
	if err != nil{
		return 0, err
	}
	if len(buf) == 0{
		return -1, nil
	}
	return strconv.Atoi(string(buf))
}

func setPid(path string, pid int) error{
	return file.WriteFile(path, pid)
}

func isRun(pid int) bool{
	if pid == -1{
		return false
	}
	c,err := getCmd(pid)
	return err == nil && len(c) > 0
}

func sendSignal(pid int,msg syscall.Signal) error{
	p,err := os.FindProcess(pid)
	if err != nil{
		return err
	}
	return p.Signal(msg)
}

func getCmd(pid int) (string,error){
	buf, err := file.ReadFile(fmt.Sprintf("/proc/%d/cmdline",pid))
	if err != nil{
		return "", err
	}
	return string(buf), err
}

func reload(pid int) error{
	c,err := getCmd(pid)
	if err != nil{
		return err
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
	err = sendSignal(pid,syscall.SIGTERM)
	if err != nil{
		return err
	}
	cmd := exec.Command(args[0],args[1:]...)
	return cmd.Start()
}

func stop(pid int) error{
	return sendSignal(pid, syscall.SIGTERM)
}

func isRoot() bool{
	return os.Geteuid() == 0
}