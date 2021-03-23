package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var log = logrus.WithFields(logrus.Fields{"package": "config"})

func New() *Config {
	return &Config{}
}

type Config struct {
	ConfigInfo *ini.File
}

//LoadConfig 加载配置文件
func (con *Config) LoadConfig(fileName string) {
	cfg, err := ini.Load(fileName)
	if err != nil {
		log.Panic("加载配置文件出错")
	}
	con.ConfigInfo = cfg
	log.Info("LoadConfig succ")
}
