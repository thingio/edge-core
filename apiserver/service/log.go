package service

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/proto/stderr"
)

func AddLogWS(ws *restful.WebService) *restful.WebService {

	for _, kind := range []*resource.Kind{resource.KindPipeTask, resource.KindFunclet, resource.KindApplet} {
		api := &logApi{kind}
		AddLog(ws, kind, api.GetLog)
	}

	return ws
}

func AddLog(ws *restful.WebService, kind *resource.Kind, function restful.RouteFunction) *restful.WebService {
	return ws.Route(ws.GET(fmt.Sprintf("/%s/{id}", kind.Name)).To(function).
		Metadata(restfulspec.KeyOpenAPITags, []string{kind.Name}).
		Doc(fmt.Sprintf("get logs of %s", kind.Name)).
		Param(ws.PathParameter("id", fmt.Sprintf("%s id", kind.Name)).DataType("string")).
		Param(ws.QueryParameter("offset", "log start offset").DataType("int")).
		Param(ws.QueryParameter("size", "log fetch size").DataType("int")).
		Writes(LogData{}))
}

type logApi struct {
	Kind *resource.Kind
}

type LogData struct {
	Content string `json:"content,omitempty"`
	Offset  string `json:"offset,omitempty"`
	Size    string `json:"size,omitempty"`
	Total   string `json:"total,omitempty"`
}

func (this *logApi) GetLog(request *restful.Request, response *restful.Response) {
	WriteError(response, stderr.NotImplemented.Of("log api of %s not implemented yet", this.Kind.Name))
}
