package cmd

import (
	"flag"
	"fmt"
	"os"
	"viv/controller"
	"viv/local/config"
	"viv/local/file"
	gateway "viv/network/gateway/tcp"
	"viv/network/tcp"
)

var(
	f	string
	s	string
	u	string
	h	bool
	v	bool
	i	bool
	d	bool
)

var version = "V4.0.0"

func init(){
	flag.StringVar(&f, "f", "F://project//golang//qLite//V4//sev//config.xml", "config file path")
	flag.StringVar(&s, "s", "", "send signal to qlite process: [stop,reload]")
	flag.StringVar(&u, "u", "http://127.0.0.1","set config download url in init")
	flag.BoolVar(&h, "h", false, "show help")
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&i, "i", false, "`init` qlite server")
	flag.BoolVar(&d, "d", false, "to `daemon`")
}

func Cmd(){
	flag.Parse()
	showInfo()
	if !isRoot(){
		fmt.Println("please run qlite as root")
		os.Exit(1)
	}
	initQLite()
	signal()
	if isQLiteRun(){
		fmt.Println("qlite is running")
		os.Exit(1)
	}
	daemon()
	if err := start(); err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func showInfo(){
	ok := h || v
	if h{
		flag.PrintDefaults()
	}
	if v{
		fmt.Println(version)
	}
	if ok{
		os.Exit(0)
	}
}

func initQLite(){
	if !i{
		return
	}
	if u[len(u) - 1] != '/'{
		u += "/"
	}
	u += "qlite/" + version + "/" + "config.xml"
	err := file.CreateDir("/etc/qlite")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.CreateDir("/etc/qlite/plugins")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.CreateDir("/etc/qlite/persistence")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.CreateFile("/etc/qlite/qlite.pid")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.DownloadFile("/etc/qlite/config.xml", u)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("success init /etc/qlite")
	os.Exit(0)
}

func signal(){
	switch s {
	case "stop":
		pid, err := getPid("/etc/qlite/qlite.pid")
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
		err = stop(pid)
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
	case "reload":
		pid, err := getPid("/etc/qlite/qlite.pid")
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
		err = reload(pid)
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		return
	}
	os.Exit(0)
}

func daemon(){
	if !isChildProcess() && d{
		if err := toDaemon();err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func start() error{
	err := setPid("/etc/qlite/qlite.pid", os.Getpid())
	if err != nil{
		return err
	}
	c, err := config.Read(f)
	if err != nil{
		return err
	}
	err = controller.Init(c)
	if err != nil{
		return err
	}
	m, err := tcp.New(c.IP, c.Port)
	if err != nil{
		return err
	}
	go kill(m)
	return m.Listen(gateway.New)
}

func Start(c config.Config) error{
	err := controller.Init(c)
	if err != nil{
		return err
	}
	m, err := tcp.New(c.IP, c.Port)
	if err != nil{
		return err
	}
	return m.Listen(gateway.New)
}

func isQLiteRun() bool{
	pid, err := getPid("/etc/qlite/qlite.pid")
	if err != nil{
		fmt.Println("/etc/qlite/qlite.pid get pid :", err)
		os.Exit(1)
	}
	if pid == -1{
		return false
	}
	return isRun(pid)
}