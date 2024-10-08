package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/conf/model"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tails"
	"strings"
)

// 指定目录下的日志文件，发送到Kafka中
func main() {
	// 1. 读配置
	iPv4, err := common.GetLocalIPv4()
	if err != nil {
		logrus.Errorf("get local ipv4 fail err:%v", err)
	}

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
	err = kafka.InitKafka(strings.Split(configObj.KafkaConfig.Address, ","), configObj.ChainSize)
	if err != nil {
		logrus.Errorf("Fail to connect kafka, err: %v", err)
		return
	}

	logrus.Infof("kafka connect success")

	// 从etcd统一获取配置项
	// 初始化etcd
	err = etcd.Init(strings.Split(configObj.EtcdConfig.Address, ","))
	if err != nil {
		logrus.Errorf("Fail to connect etcd, err: %v", err)
		return
	}

	// 获取配置
	confs, err := etcd.GetCollectorConf(configObj.GetCollectKey(iPv4))
	if err != nil {
		logrus.Errorf("Fail to get tails confs, err: %v", err)
		return
	}

	// 监控对应key的变化
	go etcd.Watch(configObj.GetCollectKey(iPv4))

	// 3. 通过tail将日志读取到内存
	err = tails.Init(confs)
	if err != nil {
		logrus.Errorf("Fail to init tailCollector, err: %v", err)
		return
	}
	logrus.Infof("tailCollector init success")

	select {}
}
