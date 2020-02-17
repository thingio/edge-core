package main

import (
	"github.com/emicklei/go-restful"
	"github.com/thingio/edge-core/common/toolkit"
	"github.com/thingio/edge-core/gateway/service"
	"log"
	"net/http"
)

type Config struct {
	NodeId string `yaml:"node_id"`
	Server struct {
		Addr   string `yaml:"addr"`
	} `yaml:"server"`
}

var config = Config{}

func init() {
	toolkit.LoadConfig(&config, "etc/edge.yaml")
}

func main() {
	restful.Add(service.NewEdgeSwaggerAPI("/apidocs"))
	err := http.ListenAndServe(config.Server.Addr, nil)
	log.Panic("Failed to start apiserver server:", err)
}
