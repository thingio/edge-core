package main

import (
	"github.com/emicklei/go-restful"
	"github.com/thingio/edge-core/apiserver/conf"
	"github.com/thingio/edge-core/apiserver/service"
	"github.com/thingio/edge-core/common/log"
	"net/http"
)

func init() {
	conf.Load("etc/apiserver.yaml")
	log.Init(conf.Config.Log)
}

func main() {
	MountAPI()
	err := http.ListenAndServe(conf.Config.Server.Addr, nil)
	log.WithError(err).Fatal("Failed to start apiserver server")
}

var API_ROOT = "/api/v1"
func MountAPI() {
	restful.Add(service.NewResourceAPI(API_ROOT))
	restful.Add(service.NewEdgeSwaggerAPI("/apidocs"))
}