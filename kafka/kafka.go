package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var Client sarama.SyncProducer

func CollectKafka(address []string) (err error) {
	// 1 生产者配置
	config := sarama.NewConfig()
	// 设置ack的确认方式
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 设置发送消息的分区，随机发送
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的信息
	config.Producer.Return.Successes = true

	// 2 连接Kafka
	Client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Errorf("producer create fail, err:%v", err)
		return err
	}

	return nil
}
