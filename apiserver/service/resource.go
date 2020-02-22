package service

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/apiserver/conf"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/datahub/api"
	"net/http"
)

func AddNodeWebService(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	api := &nodeAPI{cli}
	apiTags := []string{string(resource.KindNode)}

	ws.Route(ws.GET("/node").To(api.GetNode).
		Doc("get current node info").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Writes(resource.Node{}))
	return ws
}


type nodeAPI struct {
	Client api.DatahubApi
}

func (this *nodeAPI) GetNode(request *restful.Request, response *restful.Response) {
	res, err := this.Client.GetResource(resource.KindNode, conf.Config.NodeId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, res)
}

func AddResourceWebService(kind resource.Kind, cli api.DatahubApi, ws *restful.WebService) *restful.WebService {

	crudApi := &crudResourceAPI{kind, cli}

	singleUrl := fmt.Sprintf("/%s/{id}", kind)
	pluralUrl := fmt.Sprintf("/%s", kind)

	resourceSample := kind.NewEmptyResource()
	keySample := new(resource.ResourceKey)

	apiTags := []string{string(kind)}

	ws.Route(ws.GET(singleUrl).To(crudApi.Get).
		Doc(fmt.Sprintf("get %s by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id", string(kind)+" id").DataType("string")).
		Writes(resourceSample))

	ws.Route(ws.GET(pluralUrl).To(crudApi.List).
		Doc(fmt.Sprintf("list %s", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Reads(resourceSample).
		Writes(keySample))

	ws.Route(ws.POST(pluralUrl).To(crudApi.Create).
		Doc(fmt.Sprintf("create %s", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Reads(resourceSample).
		Writes(keySample))

	ws.Route(ws.POST(singleUrl).To(crudApi.Update).
		Doc(fmt.Sprintf("update %s by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id", string(kind)+" id").DataType("string")).
		Reads(resourceSample).
		Writes(keySample))

	ws.Route(ws.DELETE(singleUrl).To(crudApi.Delete).
		Doc(fmt.Sprintf("delete %s by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id", string(kind)+" id").DataType("string")).
		Writes(keySample))

	return ws
}

type crudResourceAPI struct {
	Kind   resource.Kind
	Client api.DatahubApi
}

func (this *crudResourceAPI) Get(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	r, err := this.Client.GetResource(this.Kind, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r)
}

func (this *crudResourceAPI) List(request *restful.Request, response *restful.Response) {
	r, err := this.Client.ListResources(this.Kind)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r)
}

func (this *crudResourceAPI) Create(request *restful.Request, response *restful.Response) {
	r := this.Kind.NewResource()
	if err := request.ReadEntity(r.Value); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	r.ResourceKey.NodeId = conf.Config.NodeId
	err := this.Client.SaveResource(r)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r.ResourceKey)
}

func (this *crudResourceAPI) Update(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	r := this.Kind.NewResourceWithId(id)
	if err := request.ReadEntity(r.Value); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	r.ResourceKey.NodeId = conf.Config.NodeId
	err := this.Client.SaveResource(r)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r.ResourceKey)
}

func (this *crudResourceAPI) Delete(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	err := this.Client.DeleteResource(this.Kind, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, resource.ResourceKey{conf.Config.NodeId, this.Kind, id})
}
