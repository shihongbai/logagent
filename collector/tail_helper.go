package collector

import (
	"errors"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var tailCollector *tail.Tail

func Init(filePath string) (err error) {
	if filePath == "" {
		e := errors.New("log filePath is empty")
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

	tailCollector, err = tail.TailFile(filePath, config)
	if err != nil {
		logrus.Errorf("create tailObj fail err:%v", err)
		return err
	}

	return nil
}

func GetTailLiens() <-chan *tail.Line {
	return tailCollector.Lines
}
