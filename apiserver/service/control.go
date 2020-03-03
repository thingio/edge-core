package service

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/proto/stderr"
	"github.com/thingio/edge-core/datahub/api"
)

func AddControlWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {

	// start/stop actions for pipetask funclet and applet
	for _, kind := range []*resource.Kind{resource.KindPipeTask, resource.KindFunclet, resource.KindApplet} {
		api := &statusControlApi{cli, kind}
		AddAction(ws, kind, "POST", "start", api.Start)
		AddAction(ws, kind, "POST", "stop", api.Stop)
	}

	// mm device control
	devApi := &deviceControlApi{}
	AddAction(ws, resource.KindDevice, "GET", "preview", devApi.Preview)
	AddAction(ws, resource.KindDevice, "GET", "snapshot", devApi.Snapshot)

	// applets/funclets call
	callApi := &callServiceApi{}
	AddAction(ws, resource.KindFunclet, "POST", "call", callApi.CallFunc)
	AddAction(ws, resource.KindApplet, "POST", "call", callApi.CallApp, "path")

	return ws
}

func AddAction(ws *restful.WebService, kind *resource.Kind, httpMethod string, action string, function restful.RouteFunction, queries ...string) *restful.WebService {

	builder := ws.Method(httpMethod).Path(fmt.Sprintf("/%s/{id}/%s", kind.Name, action)).To(function).
		Metadata(restfulspec.KeyOpenAPITags, []string{kind.Name}).
		Doc(fmt.Sprintf("%s %s", action, kind.Name)).
		Param(ws.PathParameter("id", fmt.Sprintf("%s id", kind.Name)).DataType("string")).
		Writes(resource.ResourceState{State: &resource.State{}})

	for _, q := range queries {
		builder.Param(ws.QueryParameter(fmt.Sprintf("%s={%s}", q, q), q))
	}

	return ws.Route(builder)
}

type callServiceApi struct {
}

func (this *callServiceApi) CallFunc(request *restful.Request, response *restful.Response) {
	//TODO
	WriteResult(response, stderr.NotImplemented.Of("call func api"))
}

func (this *callServiceApi) CallApp(request *restful.Request, response *restful.Response) {
	//TODO
	//path := request.QueryParameter("path")
	WriteResult(response, stderr.NotImplemented.Of("call app api"))
}

type deviceControlApi struct {
}

func (this *deviceControlApi) Preview(request *restful.Request, response *restful.Response) {
	//TODO
	WriteResult(response, stderr.NotImplemented.Of("mm device preview api"))
}

func (this *deviceControlApi) Snapshot(request *restful.Request, response *restful.Response) {
	//TODO
	WriteResult(response, stderr.NotImplemented.Of("mm device snapshot api"))
}

type statusControlApi struct {
	Client api.DatahubApi
	Kind   *resource.Kind
}

func (this *statusControlApi) Start(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	r, s, err := this.getResourceAndState(id)
	if err != nil {
		WriteError(response, stderr.ResourceAccessFailure.Error(err))
		return
	}

	state := s.Value.(*resource.State)

	currStatus := (*state)[resource.StateKeyStatus]
	if currStatus == resource.StatusInit || currStatus == resource.StatusUp {
		WriteError(response, stderr.InvalidStatusTransition.Of("%s is already in Init/Up status", r.Id))
		return
	}

	(*state)[resource.StateKeyStatus] = resource.StatusInit
	this.Client.SaveResource(s)
	WriteResult(response, resource.MakeResourceState(r, s))
}

func (this *statusControlApi) Stop(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	r, s, err := this.getResourceAndState(id)
	if err != nil {
		WriteError(response, stderr.ResourceAccessFailure.Error(err))
		return
	}

	state := s.Value.(*resource.State)

	currStatus := (*state)[resource.StateKeyStatus]
	if currStatus == resource.StatusKill || currStatus == resource.StatusDown {
		WriteError(response, stderr.InvalidStatusTransition.Of("%s is already in Init/Up status", r.Id))
		return
	}

	(*state)[resource.StateKeyStatus] = resource.StatusKill
	this.Client.SaveResource(s)
	WriteResult(response, resource.MakeResourceState(r, s))
}

func (this *statusControlApi) getResourceAndState(id string) (*resource.Resource, *resource.Resource, error) {
	r, err := this.Client.GetResource(this.Kind, id)
	if err != nil {
		return nil, nil, err
	}

	s, err := this.Client.GetResource(resource.KindState, id)
	if err != nil {
		if stderr.IsNotFound(err) {
			s = resource.KindState.NewResourceWithId(id)
		} else {
			return nil, nil, err
		}
	}
	return r, s, nil
}
