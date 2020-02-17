module github.com/thingio/edge-core/gateway

go 1.13

replace github.com/thingio/edge-core/common => ../common

require (
	github.com/emicklei/go-restful v2.11.2+incompatible
	github.com/emicklei/go-restful-openapi v1.2.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/thingio/edge-core/common v0.0.0
)
