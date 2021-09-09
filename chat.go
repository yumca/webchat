package main

import (
	library "chat/library"
	"chat/server"
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

type err error

var (
	ctx = context.Background()
)

func main() {
	conf, err := library.GetConf()
	if err != nil {
		log.Fatal("Server GetConfig Error:", err)
	}
	//flag.Parse()
	//logFileName := flag.String("log", conf.Setting.LogFile, "日志文件路径和名称")
	logFile, err := os.OpenFile(conf.Setting.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Server SetLogFile Error:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	//获取执行参数并判断
	Cprocess := false
	Args := make([]string, len(os.Args))
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		//是否为子进程
		case "-cprocess":
			Cprocess = true
		}
	}
	//Cprocess为false则表示为父进程  判断是否需要开启后台运行
	if !Cprocess && conf.Setting.Daemonize == 1 {
		Args = append(Args, "-cprocess")
		// 将命令行参数中执行文件路径转换成可用路径
		filePath, _ := filepath.Abs(os.Args[0])
		cmd := exec.Command(filePath, Args...)
		// 将其他命令传入生成出的进程
		cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start() // 开始执行新进程，不等待新进程退出
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	serv, err := server.NewNetService(conf)
	if err != nil {
		log.Fatal("Server New Error:", err)
	}
	serv.Start()
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs
	// serv.ExitServ()
	// go func() {
	// 	for s := range sigs {
	// 		switch s {
	// 		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	// 			serv.ExitServ()
	// 		}
	// 	}
	// }()
	// for {
	// 	time.Sleep(time.Second)
	// }
	// http.Handle("/", websocket.Handler(Echo))
	//
	//if err := http.ListenAndServe(":1234", nil); err != nil {
	//   log.Fatal("ListenAndServe:", err)
	//}
}
