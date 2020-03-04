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
	"strings"
	"time"
)

func init() {
	conf.Load("etc/apiserver.yaml")
	log.Init(conf.Config.Log)
}

func main() {
	cli, err := client.NewDatahubClient(conf.Config.Mqtt, common_service.ApiServer, conf.Config.NodeId)
	if err != nil {
		log.WithError(err).Fatal("Failed to start datahub client")
		return
	}

	MountRestAPI(cli)
	MountWebSocketAPI(cli)

	err = http.ListenAndServe(conf.Config.Server.Addr, nil)
	log.WithError(err).Fatal("Failed to start apiserver server")
}

var API_ROOT = "/api/v1"

func MountRestAPI(cli api.DatahubApi) {
	restful.Add(service.NewResourceAPI(API_ROOT+"/res", cli))
	restful.Add(service.NewControlAPI(API_ROOT+"/ctl", cli))
	restful.Add(service.NewDataAPI(API_ROOT+"/data", cli))
	restful.Add(service.NewLogAPI(API_ROOT+"/log", cli))
	restful.Add(service.NewEdgeSwaggerAPI("/apidocs"))

	for _, ws := range restful.RegisteredWebServices() {
		ws.Filter(NCSACommonLogFormatLogger)
	}
}

func MountWebSocketAPI(cli api.DatahubApi) {
	service.WsApi.Init(cli)
	http.Handle("/ws/", service.WsApi.Handler())
}

var NCSACommonLogFormatLogger = func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// add Filter to all web service to intercept incoming request and log time

	var username = "-"
	if req.Request.URL.User != nil {
		if name := req.Request.URL.User.Username(); name != "" {
			username = name
		}
	}

	beginTime := time.Now()
	chain.ProcessFilter(req, resp)

	log.Infof("%s %s %s %s %s %d %d %v",
		strings.Split(req.Request.RemoteAddr, ":")[0],
		username,
		req.Request.Method,
		req.Request.URL.RequestURI(),
		req.Request.Proto,
		resp.StatusCode(),
		resp.ContentLength(),
		time.Now().Sub(beginTime),
	)

}

