package service

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/juju/errors"
	"github.com/thingio/edge-core/apiserver/conf"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/datahub/api"
	"net/http"
)

func AddNodeWebService(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	api := &nodeAPI{cli}
	apiTags := []string{resource.KindNode.Name}
	sample := resource.KindNode.SampleObject

	ws.Route(ws.GET("/node").To(api.GetNode).
		Doc("get current node info").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Returns(200, "OK", sample))

	ws.Route(ws.GET("/node_state").To(api.GetNodeState).
		Doc("get current node states").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Returns(200, "OK", sample))
	return ws
}

type nodeAPI struct {
	Client api.DatahubApi
}

func (this *nodeAPI) GetNode(request *restful.Request, response *restful.Response) {
	r, err := this.Client.GetResource(resource.KindNode, conf.Config.NodeId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r)
}

func (this *nodeAPI) GetNodeState(request *restful.Request, response *restful.Response) {
	r, err := this.Client.GetResource(resource.KindNode, conf.Config.NodeId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	s, err := this.Client.GetResource(resource.KindState, conf.Config.NodeId)
	if err != nil && !errors.IsNotFound(err) {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, resource.MakeResourceState(r,s))
}

func AddResourceWebService(kind *resource.Kind, cli api.DatahubApi, ws *restful.WebService) *restful.WebService {

	crudApi := &crudResourceAPI{kind, cli}

	url := fmt.Sprintf("/%ss", kind)
	idUrl := fmt.Sprintf("/%ss/{id}", kind)

	sample := kind.SampleObject
	sampleR := kind.NewEmptyResource()
	sampleRs := []*resource.Resource{sampleR}
	keySample := new(resource.Key)

	apiTags := []string{kind.Name}

	ws.Route(ws.GET(idUrl).To(crudApi.Get).
		Doc(fmt.Sprintf("get %s by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id", kind.Name+" id").DataType("string")).
		Returns(200, "OK", sampleR))

	ws.Route(ws.GET(url).To(crudApi.List).
		Doc(fmt.Sprintf("list %s", kind.Name)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Writes(sampleRs))

	ws.Route(ws.POST(url).To(crudApi.Create).
		Doc(fmt.Sprintf("create %s", kind.Name)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Reads(sample).
		Writes(keySample))

	ws.Route(ws.POST(idUrl).To(crudApi.Update).
		Doc(fmt.Sprintf("update %s by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id", kind.Name+" id").DataType("string")).
		Reads(sample).
		Writes(keySample))

	ws.Route(ws.DELETE(idUrl).To(crudApi.Delete).
		Doc(fmt.Sprintf("delete %s by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Param(ws.PathParameter("id", kind.Name+" id").DataType("string")).
		Writes(keySample))

	if kind.Stateful {

		stateUrl := fmt.Sprintf("/%s_states", kind)
		idStateUrl := fmt.Sprintf("/%s_states/{id}", kind)

		ws.Route(ws.GET(idStateUrl).To(crudApi.GetState).
			Doc(fmt.Sprintf("get %s and its state by id", kind)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
			Param(ws.PathParameter("id", kind.Name+" id").DataType("string")).
			Returns(200, "OK", sampleR))

		ws.Route(ws.GET(stateUrl).To(crudApi.ListState).
			Doc(fmt.Sprintf("list %s and its state", kind.Name)).Metadata(restfulspec.KeyOpenAPITags, apiTags).
			Writes(sampleRs))
	}

	return ws
}

type crudResourceAPI struct {
	Kind   *resource.Kind
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
	rs, err := this.Client.ListResources(this.Kind)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, rs)
}

func (this *crudResourceAPI) Create(request *restful.Request, response *restful.Response) {
	r := this.Kind.NewResource()
	if err := request.ReadEntity(r.Value); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	r.Key.NodeId = conf.Config.NodeId
	err := this.Client.SaveResource(r)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r.Key)
}

func (this *crudResourceAPI) Update(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	r := this.Kind.NewResourceWithId(id)
	if err := request.ReadEntity(r.Value); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	r.Key.NodeId = conf.Config.NodeId
	err := this.Client.SaveResource(r)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, r.Key)
}

func (this *crudResourceAPI) Delete(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	err := this.Client.DeleteResource(this.Kind, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, resource.Key{conf.Config.NodeId, this.Kind.Name, id})
}

func (this *crudResourceAPI) GetState(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	r, err := this.Client.GetResource(this.Kind, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	s, err := this.Client.GetResource(resource.KindState, id)
	if err != nil && !errors.IsNotFound(err) {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, resource.MakeResourceState(r, s))
}

func (this *crudResourceAPI) ListState(request *restful.Request, response *restful.Response) {
	rs, err := this.Client.ListResources(this.Kind)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	ss, err := this.Client.QueryResources(resource.KindState, this.Kind.Name)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	smap := make(map[string]*resource.Resource)
	for _, s := range ss {
		smap[s.Id] = s
	}

	result := make([]*resource.ResourceState, 0)
	for _, r := range rs {
		result = append(result, resource.MakeResourceState(r, smap[r.Id]))
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
