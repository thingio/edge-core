package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/thingio/edge-core/common/proto/data"
	"github.com/thingio/edge-core/common/proto/resource"
)

func AddDeviceDataWS(ws *restful.WebService) *restful.WebService {
	ws.Route(ws.GET("/devices/{id}").
		Metadata(restfulspec.KeyOpenAPITags, []string{resource.KindDevice.Name}).
		Doc("get time series devices data").
		Param(ws.PathParameter("id", "device id").DataType("string")).
		Param(ws.QueryParameter("limit", "size limit, for eg.: 1000").DataType("int")).
		Param(ws.QueryParameter("begin_ts", "time filter: start timestamp (ns)").DataType("long")).
		Param(ws.QueryParameter("end_ts", "time filter: begin timestamp (ns)").DataType("long")).
		Param(ws.QueryParameter("last", "fast time filter: format eg.: 1s/5m/12h").DataType("string")).
		To(func(request *restful.Request, response *restful.Response) {

			//id 		:= request.PathParameter("id")
			//last 	:= request.QueryParameter("last")
			//limit 	:= request.QueryParameter("limit")

			//TODO: please talk to tsdb to get data
			result := data.TsDataSet{}

			WriteResult(response, result)
		}).
		Returns(200, "OK", []*resource.Resource{resource.KindDevice.NewResource()}))
	return ws
}


func AddAlertDataWS(ws *restful.WebService) *restful.WebService {
	ws.Route(ws.GET("/alerts/{id}").
		Metadata(restfulspec.KeyOpenAPITags, []string{"alert"}).
		Doc("get time series devices data").
		To(func(request *restful.Request, response *restful.Response) {

			//TODO: please talk to tsdb to get alerts, attach image/data
			result := data.Alert{}

			WriteResult(response, result)
		}).
		Returns(200, "OK", []*resource.Resource{resource.KindDevice.NewResource()}))

	ws.Route(ws.GET("/alerts").
		Metadata(restfulspec.KeyOpenAPITags, []string{"alert"}).
		Doc("get time series devices data").
		Param(ws.QueryParameter("level", "level, format eg. FATAL/ERROR/WARN/INFO").DataType("string")).
		Param(ws.QueryParameter("limit", "size limit, for eg. 1000").DataType("int")).
		Param(ws.QueryParameter("begin_ts", "time filter: start timestamp (ns)").DataType("long")).
		Param(ws.QueryParameter("end_ts", "time filter: begin timestamp (ns)").DataType("long")).
		Param(ws.QueryParameter("last", "fast time filter: format eg. 1s/5m/12h").DataType("string")).
		Param(ws.QueryParameter("image", "whether to show image field").DataType("bool").DefaultValue("false")).
		Param(ws.QueryParameter("data", "whether to show data field").DataType("bool").DefaultValue("false")).
		To(func(request *restful.Request, response *restful.Response) {

			//TODO: please talk to tsdb to get alerts, do NOT attach image/data by default, to reduce payload size
			result := data.Alert{}

			WriteResult(response, result)
		}).
		Returns(200, "OK", []*resource.Resource{resource.KindDevice.NewResource()}))

	return ws
}
