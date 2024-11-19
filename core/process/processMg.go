package process

import (
	"dtc-service/core/reg"
	"log"
	"os/exec"
)

func InstallProgram(path string, param string) {
	err := StartProcessAsCurrentUser(path, param, "", true)
	if err != nil {
		log.Println("Error:安装软件失败", path, err)
	} else {
		log.Println("安装软件成功", path)
	}
}

func RunProgram(name reg.AppName) {
	info := reg.Reg.GetRegistryInfo(name)
	err := StartProcessAsCurrentUser(info.ExePath, "", "", false)
	if err != nil {
		log.Println("Error:创建进程失败", name, err)
	} else {
		log.Println("创建进程成功", name)
	}
}

func KillProgram(name reg.AppName) {
	cmd := exec.Command("taskkill", "/F", "/IM", name)
	err := cmd.Run()
	if err != nil {
		log.Println("结束进程异常:", name, err)
	}
	log.Println("结束进程成功：:", name)
}
