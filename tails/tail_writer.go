package tails

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"strings"
	"time"
)

type collector struct {
	path     string
	topic    string
	operator *tail.Tail
	ctx      context.Context
	cancel   context.CancelFunc
}

func newCollector(path string, topic string) *collector {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	tt := collector{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	file, e := tail.TailFile(path, config)
	if e != nil {
		logrus.Errorf("create tail err:%v", e)
		return nil
	}

	tt.operator = file
	return &tt
}

func (t *collector) run() {
	for {
		select {
		case <-t.ctx.Done():
			logrus.Infof("tail stop %v", t.path)
			return
		case line, ok := <-t.operator.Lines:
			if len(line.Text) == 0 || len(strings.TrimSpace(line.Text)) == 0 {
				logrus.Infof("log line is empty")
				continue
			}

			if len(strings.Trim(line.Text, "\r")) == 0 {
				logrus.Infof("log line is empty")
				continue
			}

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
			kafka.SendToChan(msg)
		}
	}
}
