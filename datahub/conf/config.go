package conf

import "github.com/thingio/edge-core/common/conf"

var Config = struct {
	NodeId string          `yaml:"node_id"`
	Mqtt   conf.MqttConfig `yaml:"mqtt"`
}{}

func Load(file string) {
	conf.LoadConfig(&Config, file)
}
