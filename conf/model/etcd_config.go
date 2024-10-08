package model

import "fmt"

type EtcdConfig struct {
	Address string `ini:"address"`
	// 每个服务器会根据自己的ip到etcd拉取对应自己的配置
	CollectKey string `ini:"collect_key"`
}

func (e *EtcdConfig) GetCollectKey(ip string) string {
	return fmt.Sprintf(e.CollectKey, ip)
}
