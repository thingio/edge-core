package conf

import (
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
)

type DatahubConfig struct {
	NodeId string          `json:"node_id" yaml:"node_id"`
	Mqtt   conf.MqttConfig `json:"mqtt" yaml:"mqtt"`
	DB     conf.DBConfig   `json:"db" yaml:"db"`
	Log    log.Config      `json:"log" yaml:"log"`
}

var Config DatahubConfig

func Load(file string) {
	conf.LoadConfig(&Config, file)
}
