package service

import (
	"github.com/emicklei/go-restful"
	"github.com/thingio/edge-core/apiserver/conf"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/datahub/client"
)

func NewResourceAPI(root string) *restful.WebService {
	cli, err := client.NewDatahubClient(conf.Config.Mqtt, service.ApiServer, conf.Config.NodeId)
	if err != nil {
		log.WithError(err).Fatal("Failed to start datahub client")
	}

	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	// Add resource CRUD API
	ws = AddNodeWebService(cli, ws)

	kinds := []resource.Kind{resource.KindDevice,
		resource.KindPipeline, resource.KindPipeTask,
		resource.KindApplet, resource.KindFunclet, resource.KindServlet}

	for _, k := range kinds {
		ws = AddResourceWebService(k, cli, ws)
	}


	// Add resource Control API
	AddPipeTaskWebService(cli, ws)

	return ws
}
