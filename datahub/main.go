package main

import (
	"github.com/thingio/edge-core/datahub/conf"
	"github.com/thingio/edge-core/datahub/resource/server"
)
func init() {
	conf.Load()
}

func main() {
	server.Start()
}
