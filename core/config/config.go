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
		log.Printf("读取配置文件失败: %v\n", err)
	}
	err = yaml.Unmarshal(data, &c.Conf)
	if err != nil {
		log.Printf("解析配置文件失败: %v\n", err)
	}
}
