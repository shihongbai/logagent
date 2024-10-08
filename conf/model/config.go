package model

type Config struct {
	KafkaConfig `ini:"kafka"`
	EtcdConfig  `ini:"etcd"`
}

type KafkaConfig struct {
	Address   string `ini:"address"`
	Topic     string `ini:"topic"`
	ChainSize int64  `ini:"chain_size"`
}
