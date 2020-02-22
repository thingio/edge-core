package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/datahub/api"
)

type pipetaskAPI struct {
	Client api.DatahubApi
}

func AddPipeTaskWebService(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	api := &pipetaskAPI{cli}
	apiTags := []string{string(resource.KindNode)}

	ws.Route(ws.POST("/task/{id}/start").To(api.StartTask).
		Doc("get pipeline").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Writes(resource.Pipeline{}).
		Returns(200, "OK", resource.Pipeline{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/task/{id}/stop").To(api.StopTask).
		Doc("get pipeline").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Writes(resource.Pipeline{}).
		Returns(200, "OK", resource.Pipeline{}).
		Returns(404, "Not Found", nil))

	return ws
}

func (this *pipetaskAPI) StartTask(request *restful.Request, response *restful.Response) {
	//TODO
}

func (this *pipetaskAPI) StopTask(request *restful.Request, response *restful.Response) {
	//TODO
}
