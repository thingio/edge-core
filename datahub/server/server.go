package server

import (
	"encoding/json"
	"github.com/juju/errors"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/resource"
	talkpb "github.com/thingio/edge-core/common/proto/talk"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/talk"
	"github.com/thingio/edge-core/datahub/conf"
	"time"
)

func Start() error {
	cli := talk.NewTClient(conf.Config.Mqtt, string(service.DataHub), talkpb.TMessageChatPrefix+"/"+conf.Config.NodeId)
	if err := cli.Connect(); err != nil {
		return err
	}
	cli.SetReqHandler(reqHandler)

	tchanKey := talkpb.TMessageDataKey(conf.Config.NodeId, resource.KindAny, "")
	if ch, err := cli.Watch(tchanKey); err != nil {
		return err
	} else {
		for tmsg := range ch {
			kind := resource.Kind(tmsg.ResourceKind())
			Save(kind, tmsg.Payload)
		}
	}

	return nil
}

var handlers = map[service.ServiceFunction]func(key resource.ResourceKey) (interface{}, error){
	service.FuncGet:  GetResource,
	service.FuncList: ListResource,
}

func GetResource(key resource.ResourceKey) (interface{}, error) {
	log.Infof("receive get call: %+v", key)
	if key.Kind == resource.KindNode {
		return &resource.Resource {
			ResourceKey: key,
			Value:       NodeData(conf.Config.NodeId),
			Ts:          time.Now().UnixNano(),
			Version:     1,
		}, nil
	}

	// TODO: should support device/pipeline/task/model/function
	return nil, errors.NotFoundf("resource %v", key)
}

func ListResource(key resource.ResourceKey) (interface{}, error) {
	return nil, errors.NotFoundf("resource %v", key)
}

func reqHandler(function string, payload []byte) (i []byte, e error) {
	f, ok := handlers[service.ServiceFunction(function)]
	if !ok {
		return nil, errors.Errorf("no handler found for %s", function)
	}
	req := resource.ResourceKey{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	rsp, err := f(req)
	if err != nil {
		return nil, err
	}

	switch v := rsp.(type) {
	case *resource.Resource:
		return resource.MarshalResource(v)
	case []*resource.Resource:
		return resource.MarshalResourceList(v)
	default:
		return json.Marshal(rsp)
	}
}

func Save(kind resource.Kind, data interface{}) {

}
