package collector

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/kafka"
	"strings"
	"time"
)

type collector struct {
	path     string
	topic    string
	operator *tail.Tail
}

func (t *collector) run() {
	for {
		line, ok := <-t.operator.Lines
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

func Init(allConfList []common.CollectEntry) (err error) {
	if len(allConfList) == 0 {
		e := errors.New("allConfList is empty")
		logrus.Error(e)
		return e
	}

	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	// 创建多个tail
	collectors := make([]collector, 0, len(allConfList))
	for _, conf := range allConfList {
		tt := collector{
			path:  conf.Path,
			topic: conf.Topic,
		}
		file, e := tail.TailFile(conf.Path, config)
		if e != nil {
			logrus.Errorf("create tail err:%v", e)
			continue
		}

		tt.operator = file

		collectors = append(collectors, tt)
	}

	if len(collectors) == 0 {
		return errors.New("create collector err, collectors is empty")
	}

	// 4. 使用saram写入到Kafka
	// 从tail --> log --> kafka client -->
	for _, collector := range collectors {
		logrus.Infof("create collector:%v", collector)
		go collector.run()
	}

	return nil
}
