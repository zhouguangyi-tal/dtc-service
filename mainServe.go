package main

import (
	"dtc-service/biz/software"
	"dtc-service/core/config"
	"dtc-service/core/net"
	"dtc-service/core/reg"
	"dtc-service/core/task"
	"github.com/kardianos/service"
	"log"
)

type MainServe struct {
	conf     config.Config
	schedule task.TaskSchedule
	wsClient net.WsClient
}

func (ms *MainServe) Init(dir string) error {
	log.Println("init daemon service")
	ms.conf.Init(dir)
	ms.schedule.Init()
	reg.Reg.Init(ms.conf.Conf.RegPath)
	ms.wsClient.Init(ms.conf.Conf.WS)
	return nil
}

func (ms *MainServe) Start(s service.Service) error {
	log.Println("start run daemon service")
	ms.schedule.Start()
	ms.wsClient.Start()

	//process.RunProgram(reg.Betterme)
	tk1 := task.Task{}
	tk1.CreateTask("安装云课堂", func() {
		software.InstallSoftware(`D:\yunketang.exe`)
	})
	ms.schedule.AddDailyTask(&tk1, 17, 6)

	return nil
}

func (ms *MainServe) Stop(s service.Service) error {
	log.Println("stop  daemon service")
	return nil
}
