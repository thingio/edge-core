package conf

import (
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
)

type ApiServerConfig struct {
	NodeId string `json:"node_id" yaml:"node_id"`
	Server struct {
		Addr   string `json:"addr" yaml:"addr"`
	} `json:"server" yaml:"server"`
	Mqtt   conf.MqttConfig `json:"mqtt" yaml:"mqtt"`
	Log    log.Config      `json:"log" yaml:"log"`
}
var Config ApiServerConfig

func Load(file string) {
	conf.LoadConfig(&Config, file)
}
