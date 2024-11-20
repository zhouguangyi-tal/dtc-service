package request

import (
	"dtc-service/core/net"
	"dtc-service/core/reg"
	"fmt"
)

func GetUpdateConfig() map[string]Version {

	param := map[string]string{
		"nowversionClc": reg.Reg.GetRegistryInfo(reg.Betterme).Version,
	}
	payLoad := make(map[string]string)
	res, err := net.PostRequest[Result]("https://test-sci-gateway.speiyou.com/config/argument/device/v1/getplan", param, payLoad)
	var config map[string]Version
	if err != nil {
		fmt.Println("zzz", res, err)
	} else {
		config = res.Result.Config
		fmt.Println("zzz", config)
	}
	return config

}
