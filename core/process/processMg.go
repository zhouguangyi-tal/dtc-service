package process

import (
	"dtc-service/core/reg"
	"log"
	"os/exec"
)

func RunProgram(name reg.AppName) {
	info := reg.Reg.GetRegistryInfo(name)
	err := CreateForegroundProcess(info.ExePath)
	if err != nil {
		log.Fatalln("Error:创建进程失败", name, err)
	} else {
		log.Println("创建进程成功", name)
	}
}

func KillProgram(name reg.AppName) {
	cmd := exec.Command("taskkill", "/F", "/IM", name)
	err := cmd.Run()
	if err != nil {
		log.Fatalln("结束进程异常:", name, err)
	}
	log.Println("结束进程成功：:", name)
}
