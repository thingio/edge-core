module github.com/thingio/edge-core/datahub

go 1.13

replace github.com/thingio/edge-core/common => ../common

require (
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/boltdb/bolt v1.3.1
	github.com/juju/errors v0.0.0-20190930114154-d42613fe1ab9
	github.com/thingio/edge-core/common v0.0.0
)
