package client

type Config struct {
	Mqtt string `yaml:"node_id"`
	Server struct {
		Addr   string `yaml:"addr"`
	} `yaml:"server"`
}