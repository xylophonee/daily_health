package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Info struct {
	Users Users `yaml:"users"`
}

type Users struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email string 	`yaml:"email"`
}

//GetConfig 获取配置数据
func GetConfig(filePath string) Info {
	config := Info{}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("解析config.yaml读取错误: %v", err)
	}

	if yaml.Unmarshal(content, &config) != nil {
		log.Fatalf("解析config.yaml出错: %v", err)
	}
	return config
}
