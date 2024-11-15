package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Conf ConfEntity
}

func (c *Config) Init(dir string) {
	c.Conf = ConfEntity{}
	data, err := os.ReadFile(dir + "/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	err = yaml.Unmarshal(data, &c.Conf)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	log.Println("zzz", c.Conf.WS)
}
