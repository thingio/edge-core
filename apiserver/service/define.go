package service

import (
	"github.com/emicklei/go-restful"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/proto/stderr"
	"github.com/thingio/edge-core/datahub/api"
	"net/http"
)

func NewResourceAPI(root string, cli api.DatahubApi) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON, restful.MIME_OCTET).Produces(restful.MIME_JSON)

	// Add resource CRUD API
	for _, k := range resource.AllKinds {
		if k == resource.KindNode {
			// Node resource only support Read API
			ws = AddNodeWS(cli, ws)
		} else {
			// Other resource will support CRUD API
			ws = AddResourceWS(k, cli, ws)
		}

		// Some resource will have list API with complicated query condition
		ws = AddResourceListWS(k, cli, ws)

		// Some resource will have summary API to get complete data in once
		ws = AddResourceSummaryWS(k, cli, ws)
	}

	return ws
}

func NewControlAPI(root string, cli api.DatahubApi) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON, restful.MIME_OCTET).Produces(restful.MIME_JSON)

	AddControlWS(cli, ws)
	return ws
}

func NewDataAPI(root string, cli api.DatahubApi) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON, restful.MIME_OCTET).Produces(restful.MIME_JSON)

	AddDeviceDataWS(ws)
	AddAlertDataWS(ws)
	return ws
}

func NewLogAPI(root string, cli api.DatahubApi) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON, restful.MIME_OCTET).Produces(restful.MIME_JSON)

	AddLogWS(ws)
	return ws
}

func WriteError(rsp *restful.Response, v stderr.StdError) {
	rsp.WriteHeaderAndEntity(http.StatusInternalServerError, v)
}

func WriteResult(rsp *restful.Response, v interface{}) {
	rsp.WriteHeaderAndEntity(http.StatusOK, v)
}
