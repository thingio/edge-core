package conf

import "github.com/thingio/edge-core/common/conf"

var Config = struct {
	NodeId string          `yaml:"node_id"`
	Mqtt   conf.MqttConfig `yaml:"mqtt"`
}{}

func Load() {
	conf.LoadConfig(&Config, "etc/edge.yaml")
}
