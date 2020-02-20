package conf

import (
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
)

var Config = struct {
	NodeId string          `yaml:"node_id"`
	Mqtt   conf.MqttConfig `yaml:"mqtt"`
	Log    log.Config      `yaml:"log"`
}{}

func Load(file string) {
	conf.LoadConfig(&Config, file)
}
