package main

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/conf/model"
	"logagent/kafka"
	"logagent/utils"
	"strings"
	"time"
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
	err = kafka.InitKafka(strings.Split(configObj.Address, ","), configObj.ChainSize)
	if err != nil {
		logrus.Errorf("Fail to connect kafka, err: %v", err)
		return
	}

	logrus.Infof("kafka connect success")

	// 3. 通过tail将日志读取到内存
	err = utils.Init(configObj.LogFilePath)
	if err != nil {
		logrus.Errorf("Fail to init tailCollector, err: %v", err)
		return
	}
	logrus.Infof("tailCollector init success")

	// 4. 使用saram写入到Kafka
	// 从tail --> log --> kafka client -->
	err = run()
	if err != nil {
		logrus.Errorf("Fail to run tailCollector, err: %v", err)
		return
	}

	return
}

// log agent业务逻辑
func run() error {

	for {
		line, ok := <-utils.TailCollector.Lines
		if !ok {
			logrus.Warnf("No log is currently written")
			time.Sleep(time.Second)
			continue
		}

		// 利用通道异步 写入Kafka
		//kafka.SendMessage(msg.Text)
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		// 传入管道中
		kafka.MsgChan <- msg
	}

	return errors.New("写日志启动失败")
}
