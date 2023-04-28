package main

import (
	"flag"
	"fmt"
	"github.com/Serendipity-sw/gutil"
	"github.com/gin-gonic/gin"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/howeyc/fsnotify"
	"github.com/swgloomy/gutil/glog"
	"os"
	"os/signal"
	"syscall"
)

var (
	rt          *gin.Engine
	pidStrPath  = "./chat-from-file.pid"
	fileDirPath = "./fileDir"
	configFn    = flag.String("config", "./config.json", "config file path") //配置文件地址
	debugFlag   = flag.Bool("d", false, "debug mode")
	webTemplate = "./web-template"
	fsWatch     *fsnotify.Watcher
)

func main() {
	flag.Parse()
	err := config.ReadCfg(*configFn)
	if err != nil {
		fmt.Printf("main ReadCfg read err! filePath: %s err: %+v \n", *configFn, err.Error())
		return
	}
	readConfig()

	serverRun(*debugFlag)
	c := make(chan os.Signal, 1)
	gutil.WritePid(pidStrPath)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	//信号等待
	<-c
	serverExit()
}

func ginInit(debug bool) {
	//设置gin的工作方式
	gin.SetMode(gutil.If(debug, gin.DebugMode, gin.ReleaseMode).(string))
	rt = gin.Default()
	//允许跨域
	rt.Use()
	setGinRouter(rt)
	go func() {
		err := rt.Run(fmt.Sprintf(":8082"))
		if err != nil {
			fmt.Printf("rt run err! err: %s \n", err.Error())
		}
	}()
}

func serverRun(debug bool) {
	gutil.LogInit(debug, "./logs")

	gutil.SetCPUUseNumber(0)
	fmt.Println("set many cpu successfully!")

	deferinit.InitAll()
	fmt.Println("init all module successfully!")

	deferinit.RunRoutines()
	fmt.Println("init all run successfully!")

	_, err := os.Stat(fileDirPath)
	dirExisted := err == nil || os.IsExist(err)
	if !dirExisted {
		err = os.Mkdir(fileDirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir run err! err: %s \n", err.Error())
		}
	}

	_, err = os.Stat(webTemplate)
	dirExisted = err == nil || os.IsExist(err)
	if !dirExisted {
		err = os.Mkdir(webTemplate, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir run err! err: %s \n", err.Error())
		}
	}

	ginInit(debug)
	fmt.Println("ginInit run successfully!")
	fileWatch()
}

func serverExit() {
	err := fsWatch.RemoveWatch(webTemplate)
	if err != nil {
		glog.Error("main dir watch remove err! webTemplate: %s err: %+v \n", webTemplate, err)
	}
	deferinit.StopRoutines()
	fmt.Println("stop routine successfully!")

	deferinit.FiniAll()
	fmt.Println("stop all modules successfully!")

	glog.Close()

	os.Exit(0)
}
