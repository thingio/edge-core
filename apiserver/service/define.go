package service

import (
	"github.com/emicklei/go-restful"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/datahub/api"
)

func NewResourceAPI(root string, cli api.DatahubApi) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	// Add resource CRUD API
	for _, k := range resource.AllKinds {
		if k == resource.KindNode {
			// Node resource only support Read API
			ws = AddNodeWS(cli, ws)
		} else {
			// Other resource will support CRUD API
			ws = AddResourceWS(k, cli, ws)
		}

		// Some resource will have summary API to get complete data in once
		ws = AddResourceSummaryWS(k, cli, ws)
	}

	return ws
}

func NewControlAPI(root string, cli api.DatahubApi) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	// Add resource Control API
	AddPipeTaskWebService(cli, ws)
	return ws
}