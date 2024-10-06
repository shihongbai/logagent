package utils

import (
	"errors"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var TailCollector *tail.Tail

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

	TailCollector, err = tail.TailFile(filePath, config)
	if err != nil {
		logrus.Errorf("create tailObj fail err:%v", err)
		return err
	}

	return nil
}
