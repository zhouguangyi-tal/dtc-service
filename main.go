package main

import (
	"fmt"
	"github.com/kardianos/service"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	SVRNAME    = "ADtcService"
	SVRDISNAME = "tal dtc server daemon"
	SVRDES     = "守护服务"
)

func main() {
	cfg := service.Config{
		Name:        SVRNAME,
		DisplayName: SVRDISNAME,
		Description: SVRDES,
	}

	//dir, err := filepath.Abs(".")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	t := time.Now()
	fn := fmt.Sprintf(dir+"\\server-%v.log", t.YearDay())
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)

	var ms = &MainServe{}
	err = ms.Init(dir) //主线初始化
	if err != nil {
		log.Printf("Init service error:%s\n", err.Error())
	}
	server := service.ChosenSystem()
	srv, err := server.New(ms, &cfg)
	if err != nil {
		log.Printf("Init service error:%s\n", err.Error())
	}
	if len(os.Args) > 1 {
		log.Println("Args[1]:", os.Args[1])
		switch os.Args[1] {
		case "install":
			err := srv.Install()
			if err != nil {
				log.Printf("Install service error:%s\n", err.Error())
			} else {
				log.Println("server is installed")
			}
			break
		case "uninstall":
			err := srv.Stop()
			err = srv.Uninstall()
			if err != nil {
				log.Printf("Uninstall service error:%s\n", err.Error())
			} else {
				log.Println("server is uninstalled")
			}
			break
		case "debug":
			err = srv.Run()
			if err != nil {
				log.Printf("debug programe error:%s\n", err.Error())
			}
			break
		default:
			log.Println("command errer")
		}
		return
	}
	err = srv.Run()
	if err != nil {
		log.Printf("run deamon error:%s\n", err.Error())
		return
	}
}
