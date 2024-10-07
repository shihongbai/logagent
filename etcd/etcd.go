package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sirupsen/logrus"
	client "go.etcd.io/etcd/clientv3"
	"logagent/common"
	"logagent/tails"
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
func GetCollectorConf(key string) (collectCfgList []*common.CollectEntry, err error) {
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

func Watch(key string) {
	for {
		// 使用可取消的上下文，确保能够在需要时停止监控
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // 在需要时手动调用 cancel 取消监控
		watcher := cli.Watch(ctx, key)

		// 监控通道中的事件
		for w := range watcher {
			// 检查 watcher 是否被取消或断开连接
			if w.Canceled {
				logrus.Errorf("watcher canceled or disconnected")
				// 根据需求处理 watcher 关闭逻辑，可能是退出或者尝试重连
				// break 或重试逻辑可以添加在这里
				return // 假设我们选择退出，避免无限循环
			}

			logrus.Infof("received new configuration from etcd")

			// 遍历每个事件
			for _, ev := range w.Events {
				// 根据事件类型处理不同的操作
				switch ev.Type {
				case mvccpb.PUT:
					logrus.Infof("new event type: %v, key: %v", ev.Type, string(ev.Kv.Key))

					var newConf []*common.CollectEntry
					// 解析新的配置数据
					err := json.Unmarshal(ev.Kv.Value, &newConf)
					if err != nil {
						logrus.Errorf("failed to unmarshal new etcd config, error: %v", err)
						continue
					}

					// 将新的配置应用到系统中
					tails.SetNewConf(newConf)

					//case mvccpb.DELETE:
					//	logrus.Infof("configuration deleted, key: %v", string(ev.Kv.Key))
					//	// 处理配置被删除的情况，例如清空配置或恢复默认
					//	tails.ClearConf() // 假设存在清空配置的函数
					//}
				}
			}

			// 其他关闭或清理资源的逻辑可在此补充
			logrus.Info("watcher exited")
		}
	}

}
