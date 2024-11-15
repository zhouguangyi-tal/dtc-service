package main

import (
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

func (ms *MainServe) Start(s service.Service) error {
	log.Println("start run daemon service")
	task1 := task.Task{}
	task1.CreatTickerTask("startPlayer", func() {
		log.Println("run task1")
		//process.RunProgram(reg.Recordplayerteaching)
	}, 10)
	ms.schedule.AddTask("startPlayer", &task1)
	ms.schedule.Start()
	ms.wsClient.Start()
	return nil
}

func (ms *MainServe) Stop(s service.Service) error {
	log.Println("stop  daemon service")
	return nil
}

func (ms *MainServe) Init(dir string) error {
	log.Println("init daemon service")
	ms.conf.Init(dir)
	ms.schedule.Init()
	reg.Reg.Init(ms.conf.Conf.RegPath)
	ms.wsClient.Init(ms.conf.Conf.WS)
	return nil
}
