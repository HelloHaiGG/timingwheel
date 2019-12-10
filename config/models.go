package config
type AppCfg struct {
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
	Hosts []string `yaml:"hosts"`
	Port  int      `yaml:"port"`
}

type Servers struct {
	Order   string `yaml:"order"`
	Finance string `yaml:"finance"`
	Gateway string `yaml:"gateway"`
}