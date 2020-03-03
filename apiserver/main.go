package main

import (
	"github.com/emicklei/go-restful"
	"github.com/thingio/edge-core/apiserver/conf"
	"github.com/thingio/edge-core/apiserver/service"
	"github.com/thingio/edge-core/common/log"
	common_service "github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/datahub/api"
	"github.com/thingio/edge-core/datahub/client"
	"net/http"
)

func init() {
	conf.Load("etc/apiserver.yaml")
	log.Init(conf.Config.Log)
}

func main() {
	cli, err := client.NewDatahubClient(conf.Config.Mqtt, common_service.ApiServer, conf.Config.NodeId)
	if err != nil {
		log.WithError(err).Fatal("Failed to start datahub client")
	}
	MountAPI(cli)

	err = http.ListenAndServe(conf.Config.Server.Addr, nil)
	log.WithError(err).Fatal("Failed to start apiserver server")
}

var API_ROOT = "/api/v1"

func MountAPI(cli api.DatahubApi) {
	restful.Add(service.NewResourceAPI(API_ROOT+"/res", cli))
	restful.Add(service.NewControlAPI(API_ROOT+"/ctl", cli))
	restful.Add(service.NewDataAPI(API_ROOT+"/data", cli))
	restful.Add(service.NewLogAPI(API_ROOT+"/log" +
		"", cli))
	restful.Add(service.NewEdgeSwaggerAPI("/apidocs"))
}
