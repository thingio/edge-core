package conf

import (
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
)

type DatahubConfig struct {
	NodeId string          `yaml:"node_id"`
	Mqtt   conf.MqttConfig `yaml:"mqtt"`
	DB     conf.DBConfig   `yaml:"db"`
	Log    log.Config      `yaml:"log"`
}

var Config DatahubConfig

func Load(file string) {
	conf.LoadConfig(&Config, file)
}
