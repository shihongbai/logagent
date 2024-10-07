package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func sendMessage() error {
	for {
		select {
		case msg := <-msgChan:
			// 发送消息
			pid, offset, err := Client.SendMessage(msg)
			if err != nil {
				logrus.Errorf("send message error:%v", err)
				continue
			}

			logrus.Infof("send msg(info: %s) to kafka success pid:%d offset:%d", msg.Value, pid, offset)
		}
	}
}

func SendToChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}

func InitKafka(address []string, chanSize int64) (err error) {
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

	// 初始化通道
	msgChan = make(chan *sarama.ProducerMessage, chanSize)

	// 启动kafka发送消息
	go sendMessage()

	return nil
}
