package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/proto/stderr"
	"github.com/thingio/edge-core/datahub/api"
)

func AddResourceListWS(kind *resource.Kind, cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	switch kind {
	case resource.KindDevice:
		ws = AddDeviceListWS(cli, ws)
	case resource.KindProduct:
		ws = AddProductListWS(cli, ws)
	case resource.KindProtocol:
		ws = AddProtocolListWS(cli, ws)
	case resource.KindPipeline:
		ws = AddPipelineListWS(cli, ws)
	case resource.KindPipeTask:
		ws = AddPipeTaskListWS(cli, ws)
		ws = AddPipeTaskStateListWS(cli, ws)
	}
	return ws
}

func AddPipeTaskStateListWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	return ws.Route(ws.GET("/pipetask_states").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindPipeTask.Name}).
		Doc("get pipetasks list").
		Param(ws.QueryParameter("genus", "filter by genus [ms|ts]").DataType("string")).
		Param(ws.QueryParameter("pipeline_id", "filter by pipeline_id").DataType("string")).
		To(
			func(request *restful.Request, response *restful.Response) {
				rs, err := cli.ListResources(resource.KindPipeTask)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}

				ss, err := cli.QueryResources(resource.KindState, resource.KindPipeTask.Name)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}

				smap := make(map[string]*resource.Resource)
				for _, s := range ss {
					smap[s.Id] = s
				}

				// filters
				genus := request.QueryParameter("genus")
				pipelineId := request.QueryParameter("pipeline_id")
				result := make([]*resource.ResourceState, 0)
				for _, r := range rs {
					dv := r.Value.(*resource.PipeTask)
					exclude := false
					exclude = exclude || (genus != "" && genus != dv.Genus)
					exclude = exclude || (pipelineId != "" && pipelineId != dv.PipelineId)
					if !exclude {
						result = append(result, resource.MakeResourceState(r, smap[r.Id]))
					}
				}
				WriteResult(response, result)
			}).
		Returns(200, "OK", []*resource.Resource{resource.KindPipeTask.NewResource()}))
}

func AddPipeTaskListWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {

	return ws.Route(ws.GET("/pipetasks").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindPipeTask.Name}).
		Doc("get pipetasks list").
		Param(ws.QueryParameter("genus", "filter by genus [ms|ts]").DataType("string")).
		Param(ws.QueryParameter("pipeline_id", "filter by pipeline_id").DataType("string")).
		To(
			func(request *restful.Request, response *restful.Response) {
				rs, err := cli.ListResources(resource.KindPipeTask)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}

				// filters
				genus := request.QueryParameter("genus")
				pipelineId := request.QueryParameter("pipeline_id")
				result := make([]*resource.Resource, 0)
				for _, r := range rs {
					dv := r.Value.(*resource.PipeTask)
					exclude := false
					exclude = exclude || (genus != "" && genus != dv.Genus)
					exclude = exclude || (pipelineId != "" && pipelineId != dv.PipelineId)
					if !exclude {
						result = append(result, r)
					}
				}
				WriteResult(response, result)
			}).
		Returns(200, "OK", []*resource.Resource{resource.KindPipeTask.NewResource()}))
}

func AddPipelineListWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {

	return ws.Route(ws.GET("/pipelines").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindPipeline.Name}).
		Doc("get pipelines list").
		Param(ws.QueryParameter("genus", "filter by genus [ms|ts]").DataType("string")).
		To(
			func(request *restful.Request, response *restful.Response) {
				rs, err := cli.ListResources(resource.KindPipeline)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}

				// filters
				genus := request.QueryParameter("genus")
				result := make([]*resource.Resource, 0)
				for _, r := range rs {
					dv := r.Value.(*resource.Pipeline)
					if !(genus != "" && genus != dv.Genus) {
						result = append(result, r)
					}
				}
				WriteResult(response, result)
			}).
		Returns(200, "OK", []*resource.Resource{resource.KindPipeline.NewResource()}))
}

func AddProtocolListWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {

	return ws.Route(ws.GET("/protocols").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindProtocol.Name}).
		Doc("get products list").
		Param(ws.QueryParameter("genus", "filter by genus [ms|ts]").DataType("string")).
		To(
			func(request *restful.Request, response *restful.Response) {
				genus := request.QueryParameter("genus")
				rs, err := cli.ListResources(resource.KindProtocol)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}
				// filters
				result := make([]*resource.Resource, 0)
				for _, r := range rs {
					dv := r.Value.(*resource.DeviceProtocol)
					if !(genus != "" && genus != dv.Genus) {
						result = append(result, r)
					}
				}
				WriteResult(response, result)
			}).
		Returns(200, "OK", []*resource.Resource{resource.KindProtocol.NewResource()}))
}

func AddProductListWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	return ws.Route(ws.GET("/products").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindProduct.Name}).
		Doc("get products list").
		Param(ws.QueryParameter("genus", "filter by genus [ms|ts]").DataType("string")).
		Param(ws.QueryParameter("protocol_id", "filter by protocol id").DataType("string")).
		To(
			func(request *restful.Request, response *restful.Response) {
				rs, err := cli.ListResources(resource.KindProduct)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}

				// filters
				genus := request.QueryParameter("genus")
				protocolId := request.QueryParameter("protocol_id")
				result := make([]*resource.Resource, 0)
				for _, r := range rs {
					dv := r.Value.(*resource.DeviceProduct)
					exclude := false
					exclude = exclude || (genus != "" && genus != dv.Genus)
					exclude = exclude || (protocolId != "" && protocolId != dv.ProtocolId)
					if !exclude {
						result = append(result, r)
					}
				}
				WriteResult(response, result)
			}).
		Returns(200, "OK", []*resource.Resource{resource.KindProduct.NewResource()}))
}

func AddDeviceListWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	return ws.Route(ws.GET("/devices").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindDevice.Name}).
		Doc("get devices list").
		Param(ws.QueryParameter("genus", "filter by genus [ms|ts]").DataType("string")).
		Param(ws.QueryParameter("product_id", "filter by product id").DataType("string")).
		To(
			func(request *restful.Request, response *restful.Response) {
				devs, err := cli.ListResources(resource.KindDevice)
				if err != nil {
					WriteError(response, stderr.ResourceAccessFailure.Error(err))
					return
				}

				// filters
				genus := request.QueryParameter("genus")
				productId := request.QueryParameter("product_id")
				result := make([]*resource.Resource, 0)
				for _, d := range devs {
					dv := d.Value.(*resource.Device)
					exclude := false
					exclude = exclude || (genus != "" && genus != dv.Genus)
					exclude = exclude || (productId != "" && productId != dv.ProductId)
					if !exclude {
						result = append(result, d)
					}
				}
				WriteResult(response, result)
			}).
		Returns(200, "OK", []*resource.Resource{resource.KindDevice.NewResource()}))
}
