package server

import (
	"github.com/juju/errors"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/talk"
	"github.com/thingio/edge-core/common/toolkit"
	"github.com/thingio/edge-core/datahub/conf"
)

func Start() error {
	cli := talk.NewTClient(conf.Config.Mqtt, string(service.DataHub), talk.TChanPrefix+"/"+conf.Config.NodeId)
	if err := cli.Connect(); err != nil {
		return err
	}
	cli.SetReqHandler(reqHandler)

	tchanKey := talk.TDataKey(conf.Config.NodeId, resource.KindAny, "")
	if ch, err := cli.Watch(tchanKey); err != nil {
		return err
	} else {
		for tmsg := range ch {
			kind := resource.Kind(tmsg.KeyPart(talk.KPI_ResourceKind))
			Save(kind, tmsg.Payload)
		}
	}

	return nil
}


var handlers = map[service.ServiceFunction]func(key resource.ResourceKey) (interface{}, error) {
	service.FuncGet: GetResource,
	service.FuncList: ListResource,
}


func GetResource(key resource.ResourceKey) (interface{}, error) {
	log.Infof("receive get call: %+v", key)
	if key.Kind == resource.KindNode {
		return toolkit.NodeData(conf.Config.NodeId), nil
	}

	// TODO: should support device/pipeline/task/model/function
	return nil, nil
}

func ListResource(key resource.ResourceKey) (interface{}, error) {
	return nil, nil
}


func reqHandler(function string, payload interface{}) (i interface{}, e error) {
	f, ok := handlers[service.ServiceFunction(function)]
	if !ok {
		return nil, errors.Errorf("no handler found for %s", function)
	}
	args, ok := payload.(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("payload is not parsable")
	}

	key := resource.ResourceKey{
		args["node_id"].(string),
		resource.Kind(args["kind"].(string)),
		args["id"].(string),
	}
	return f(key)
}

func Save(kind resource.Kind, data interface{}) {

}
