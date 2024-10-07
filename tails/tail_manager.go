package tails

import (
	"errors"
	"github.com/sirupsen/logrus"
	"logagent/common"
)

type TailManager struct {
	TailManager  map[string]*collector
	CollectEntry []*common.CollectEntry
	ConfChan     chan []*common.CollectEntry
}

func (m *TailManager) watch() {
	for {
		newConf := <-m.ConfChan
		for _, entry := range newConf {
			// tail更新
			if m.isExist(entry) {
				continue
			}

			// 原先没有的, 创建一个新的tail
			c := newCollector(entry.Path, entry.Topic)
			logrus.Infof("new collector success run, topic: %s, path: %s", entry.Topic, entry.Path)
			m.TailManager[entry.Path] = c
			go c.run()
		}

		// 找出manager有的，但是new没有的，然后停掉
		for k, v := range m.TailManager {
			var found bool
			for _, conf := range newConf {
				if k == conf.Path {
					found = true
					break
				}
			}

			if !found {
				// todo 停止时候要记录读取文件的最后一次位置
				logrus.Infof("tail manager remove path: %s", k)
				m.stopTail(v)
				m.deleteTail(k)
			}
		}
	}
}

func (m *TailManager) isExist(conf *common.CollectEntry) bool {
	_, ok := m.TailManager[conf.Path]

	return ok
}

func (m *TailManager) deleteTail(k string) {
	delete(m.TailManager, k)
}

func (m *TailManager) stopTail(v *collector) {
	v.cancel()
}

var (
	ttMgr *TailManager
)

func Init(allConfList []*common.CollectEntry) (err error) {
	if len(allConfList) == 0 {
		e := errors.New("allConfList is empty")
		logrus.Error(e)
		return e
	}

	ttMgr = &TailManager{
		TailManager:  make(map[string]*collector, 20),
		CollectEntry: allConfList,
		ConfChan:     make(chan []*common.CollectEntry),
	}

	// 创建多个tail
	collectors := make([]*collector, 0, len(allConfList))
	for _, conf := range allConfList {
		c := newCollector(conf.Path, conf.Topic)

		if c == nil {
			logrus.Error("New collector fail")
		}

		// 初始化tail管理者，用于后期tail的更新
		ttMgr.TailManager[c.path] = c
		collectors = append(collectors, c)
	}

	if len(collectors) == 0 {
		return errors.New("create tails err, collectors is empty")
	}

	// 4. 使用saram写入到Kafka
	// 从tail --> log --> kafka client -->
	for _, c := range collectors {
		logrus.Infof("create tails:%v", c)
		go c.run()
	}

	// 开启通道监听新的tail，并管理之前已经创建的配置
	ttMgr.ConfChan = make(chan []*common.CollectEntry) // 阻塞chan

	// 等待新的配置
	go ttMgr.watch()

	return nil
}

func SetNewConf(newConf []*common.CollectEntry) {
	ttMgr.ConfChan <- newConf
}
