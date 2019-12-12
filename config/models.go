package config

type appCfg struct {
	LTopic  LogTopic `yaml:"l_topic"`
	Kafka   KafkaCfg `yaml:"kafka"`
	Servers Servers  `yaml:"servers"`
}

type LogTopic struct {
	Order   string `yaml:"order"`
	Finance string `yaml:"finance"`
	Gateway string `yaml:"gateway"`
}

type KafkaCfg struct {
	Brokers []string `yaml:"brokers"`
	Timeout int      `yaml:"timeout"`
}

type Servers struct {
	Order   string `yaml:"order"`
	Finance string `yaml:"finance"`
	Gateway string `yaml:"gateway"`
}
