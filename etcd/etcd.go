package etcd

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	client "go.etcd.io/etcd/clientv3"
	"logagent/common"
	"time"
)

// 连接etcd获取远程配置项
var (
	cli *client.Client
)

func Init(address []string) error {
	c, err := client.New(client.Config{
		Endpoints:   address,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		logrus.Errorf("init etcd client fail,err:%v", err)
		return err
	}

	cli = c
	return nil
}

// GetCollectorConf 拉取日志收集的配置
func GetCollectorConf(key string) (collectCfgList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get etcd config fail,err:%v", err)
		return
	}

	if len(res.Kvs) == 0 {
		logrus.Warnf("get etcd config fail by key, key:%v", key)
		return
	}

	ret := res.Kvs[0]

	err = json.Unmarshal(ret.Value, &collectCfgList)
	if err != nil {
		logrus.Errorf("unmarshal etcd config fail,err:%v", err)
		return
	}

	return
}
