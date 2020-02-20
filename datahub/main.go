package main

import (
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/datahub/conf"
	"github.com/thingio/edge-core/datahub/server"
)
func init() {
	conf.Load("etc/datahub.yaml")
	log.Init(conf.Config.Log)
}

func main() {
	err := server.Start()
	log.WithError(err).Errorf("fail to start server")
}
