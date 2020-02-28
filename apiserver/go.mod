module github.com/thingio/edge-core/apiserver

go 1.13

replace (
	github.com/thingio/edge-core/common => ../common
	github.com/thingio/edge-core/datahub => ../datahub
)

require (
	github.com/emicklei/go-restful v2.11.2+incompatible
	github.com/emicklei/go-restful-openapi v1.2.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/juju/errors v0.0.0-20190930114154-d42613fe1ab9
	github.com/thingio/edge-core/common v0.0.0
	github.com/thingio/edge-core/datahub v0.0.0-00010101000000-000000000000
)
