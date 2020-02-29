package service

import (
	"fmt"
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
	apiTags := []string{resource.KindPipeTask.Name}

	ws.Route(ws.POST(fmt.Sprintf("/%s/{id}/start",resource.KindPipeTask.Name)).To(api.StartTask).
		Doc("start pipestart").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id","pipetask id").DataType("string")).
		Writes(resource.PipeTask{}).
		Returns(200, "OK", resource.PipeTask{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST(fmt.Sprintf("/%s/{id}/stop",resource.KindPipeTask.Name)).To(api.StopTask).
		Doc("stop pipeline").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id","pipetask id").DataType("string")).
		Writes(resource.PipeTask{}).
		Returns(200, "OK", resource.PipeTask{}).
		Returns(404, "Not Found", nil))

	return ws
}

func (this *pipetaskAPI) StartTask(request *restful.Request, response *restful.Response) {
	//TODO
}

func (this *pipetaskAPI) StopTask(request *restful.Request, response *restful.Response) {
	//TODO
}
