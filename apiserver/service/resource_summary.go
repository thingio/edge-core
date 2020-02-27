package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/datahub/api"
	"net/http"
)

func AddResourceSummaryWS(kind *resource.Kind, cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	switch kind {
	case resource.KindDevice:
		ws = AddDeviceSummaryWS(cli, ws)
	}

	return ws
}

type DeviceSummary struct {
	resource.Device
	ProductName  string `json:"product_name,omitempty"`  // auxiliary fields, updated according to product_id
	ProtocolId   string `json:"protocol_id,omitempty"`   // same as above
	ProtocolName string `json:"protocol_name,omitempty"` // same as above
}

func AddDeviceSummaryWS(cli api.DatahubApi, ws *restful.WebService) *restful.WebService {
	apiTags := []string{resource.KindDevice.Name}

	ws.Route(ws.GET("/device_summarys").To(func(request *restful.Request, response *restful.Response) {
		devs, err := cli.ListResources(resource.KindDevice)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		prods, err := cli.ListResources(resource.KindProduct)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		prodMap := make(map[string]*resource.DeviceProduct)
		for _, p := range prods {
			prodMap[p.Id] = p.Value.(*resource.DeviceProduct)
		}
		protos, err := cli.ListResources(resource.KindProtocol)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		protoMap := make(map[string]*resource.DeviceProtocol)
		for _, p := range protos {
			protoMap[p.Id] = p.Value.(*resource.DeviceProtocol)
		}

		for _, d := range devs {
			dv := d.Value.(*resource.Device)
			if prod, ok := prodMap[dv.ProductId]; ok {
				if proto, ok := protoMap[prod.ProtocolId]; ok {
					d.Value = &DeviceSummary{Device: *dv, ProductName: prod.Name, ProtocolId: proto.Id, ProtocolName: proto.Name}
				}
			}

		}
		response.WriteHeaderAndEntity(http.StatusOK, devs)
	}).
		Doc("get device summarys including product and protocol").Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Returns(200, "OK", []DeviceSummary{}))
	return ws
}
