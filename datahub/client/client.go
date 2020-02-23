package client

import (
	"encoding/json"
	"errors"
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/proto/resource"
	talkpb "github.com/thingio/edge-core/common/proto/talk"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/talk"
	"github.com/thingio/edge-core/datahub/api"
	"log"
	"time"
)

type DatahubClient struct {
	*talk.TClient
	NodeId string
}

func NewDatahubClient(cfg conf.MqttConfig, serviceId service.ServiceId, nodeId string) (api.DatahubApi, error) {
	tclient := talk.NewTClient(cfg, string(serviceId), talkpb.TMessageChatPrefix+"/"+nodeId)
	if err := tclient.Connect(); err != nil {
		return nil, err
	}
	return &DatahubClient{tclient, nodeId}, nil
}

func (dc *DatahubClient) GetResource(kind resource.Kind, id string) (*resource.Resource, error) {
	msg := talkpb.NewTMessage()
	msg.Method = talkpb.MethodREQ
	msg.Key = talkpb.TMessageChatKey(dc.NodeId, service.DataHub, service.FuncGet)
	payload, err := json.Marshal(resource.ResourceKey{dc.NodeId, kind, id})
	if err != nil {
		return nil, err
	}
	msg.Payload = payload
	rsp, err := dc.TClient.Send(msg)
	if err != nil {
		return nil, err
	}

	if rsp.Method == talkpb.MethodERR {
		return nil, errors.New(string(rsp.Payload))
	} else if rsp.Method == talkpb.MethodRSP {
		r, err := resource.UnmarshalResource(kind, rsp.Payload)
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		log.Fatal("No other talk methods should go here!")
		return nil, nil
	}
}

func (dc *DatahubClient) ListResources(kind resource.Kind) ([]*resource.Resource, error) {
	msg := talkpb.NewTMessage()
	msg.Method = talkpb.MethodREQ
	msg.Key = talkpb.TMessageChatKey(dc.NodeId, service.DataHub, service.FuncList)
	payload, err := json.Marshal(resource.ResourceKey{NodeId: dc.NodeId, Kind: kind})
	if err != nil {
		return nil, err
	}
	msg.Payload = payload
	rsp, err := dc.TClient.Send(msg)
	if err != nil {
		return nil, err
	}

	if rsp.Method == talkpb.MethodERR {
		return nil, errors.New(string(rsp.Payload))
	} else if rsp.Method == talkpb.MethodRSP {
		rs, err := resource.UnmarshalResourceList(kind, rsp.Payload)
		if err != nil {
			return nil, err
		}
		return rs, nil
	} else {
		log.Fatal("No other talk methods should go here!")
		return nil, nil
	}
}

func (dc *DatahubClient) WatchResource(kind resource.Kind, watcher api.ResourceWatcher) error {
	key := talkpb.TMessageDataKey(dc.NodeId, kind, "#")
	tmsgCh, err := dc.TClient.Watch(key)
	if err != nil {
		return err
	}

	for tmsg := range tmsgCh {
		if tmsg.Sender != dc.TClient.ClientId() { // ignore the change made by client itself
			res, err := resource.UnmarshalResource(kind, tmsg.Payload)
			if err != nil {
				return err
			}
			watcher(res)
		}
	}
	return nil
}

func (dc *DatahubClient) SaveResource(r *resource.Resource) error {
	msg := talkpb.NewTMessage()
	msg.Method = talkpb.MethodPUB
	msg.Key = talkpb.TMessageDataKey(dc.NodeId, r.Kind, r.Id)
	payload, err := resource.MarshalResource(r)
	if err != nil {
		return err
	}
	msg.Payload = payload
	_, err = dc.TClient.Send(msg)
	return err
}

func (dc *DatahubClient) DeleteResource(kind resource.Kind, id string) error {
	msg := talkpb.NewTMessage()
	msg.Method = talkpb.MethodPUB
	msg.Key = talkpb.TMessageDataKey(dc.NodeId, kind, id)
	emptyRes := &resource.Resource{
		ResourceKey: resource.ResourceKey{dc.NodeId, kind, id},
		Ts:          time.Now().UnixNano(),
		Version:     0,
	}
	payload, err := resource.MarshalResource(emptyRes)
	if err != nil {
		return err
	}
	msg.Payload = payload
	_, err = dc.TClient.Send(msg)
	return err
}
