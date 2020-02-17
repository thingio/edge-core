package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"net/http"
	"path"
)

var SwaggerUIRoot = "public/swagger-ui"

func NewEdgeSwaggerAPI(root string) *restful.WebService {
	api := new(restful.WebService)
	api.Path(root).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	api.Route(api.GET("/swagger.json").To(GetSwagger))
	api.Route(api.GET("/").To(GetSwaggerUI))
	api.Route(api.GET("/{subpath:*}").To(GetSwaggerUIResouces))
	return api
}

func GetSwagger(request *restful.Request, response *restful.Response) {
	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
	}
	swagger := restfulspec.BuildSwagger(config)
	response.WriteEntity(swagger)
}

func GetSwaggerUI(request *restful.Request, response *restful.Response) {
	http.ServeFile(
		response.ResponseWriter,
		request.Request,
		SwaggerUIRoot+"/index.html")
}

func GetSwaggerUIResouces(req *restful.Request, resp *restful.Response) {
	actual := path.Join(SwaggerUIRoot, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request, actual)
}
