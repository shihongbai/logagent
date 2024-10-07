package model

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address   string `ini:"address"`
	Topic     string `ini:"topic"`
	ChainSize int64  `ini:"chain_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}
