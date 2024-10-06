package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/conf/model"
	"logagent/kafka"
	"strings"
)

// 指定目录下的日志文件，发送到Kafka中

func main() {
	// 1. 读配置
	cfg, err := ini.Load("./conf/config.ini")
	if err != nil {
		logrus.Errorf("Fail to load config.ini, err: %v", err)
		return
	}

	configObj := new(model.Config)
	err = cfg.MapTo(configObj)
	if err != nil {
		logrus.Errorf("Fail to parse config.ini, err: %v", err)
		return
	}

	// 2. 初始化连接kafka
	err = kafka.CollectKafka(strings.Split(configObj.Address, ","))
	if err != nil {
		logrus.Errorf("Fail to connect kafka, err: %v", err)
		return
	}

	logrus.Infof("kafka connect success")
	defer kafka.Client.Close()
	// 3. 通过tail将日志读取到内存
	// 4. 使用saram写入到Kafka
}
