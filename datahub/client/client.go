package client

import (
	"errors"
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/talk"
	"github.com/thingio/edge-core/datahub/api"
	"log"
	"time"
)

type DatahubClient struct {
	*talk.TClient
	NodeId  string
}

func NewDatahubClient(cfg conf.MqttConfig, serviceId service.ServiceId, nodeId string) (api.DatahubApi, error) {
	tclient := talk.NewTClient(cfg, string(serviceId), talk.TChanPrefix+"/"+nodeId)
	if err := tclient.Connect(); err != nil {
		return nil, err
	}
	return &DatahubClient{tclient, nodeId}, nil
}

func (dc *DatahubClient) GetResource(kind resource.Kind, id string) (*resource.Resource, error) {
	msg := talk.NewTMsg()
	msg.Method = talk.REQ
	msg.Key = talk.TChanKey(dc.NodeId, service.DataHub, service.FuncList)
	msg.Payload = resource.ResourceKey{dc.NodeId, kind, id}
	rspCh, err := dc.TClient.Send(msg)
	if err != nil {
		return nil, err
	}
	defer close(rspCh)

	rsp := <-rspCh
	if rsp.Method == talk.ERR {
		return nil, errors.New(rsp.Payload.(string))
	} else if rsp.Method == talk.RSP {
		r := rsp.Payload.(*resource.Resource)
		return r, nil
	} else {
		log.Fatal("No other talk methods should go here!")
		return nil, nil
	}
}

func (dc *DatahubClient) ListResources(kind resource.Kind) ([]*resource.Resource, error) {
	msg := talk.NewTMsg()
	msg.Method = talk.REQ
	msg.Key = talk.TChanKey(dc.NodeId, service.DataHub, service.FuncList)
	msg.Payload = resource.ResourceKey{NodeId: dc.NodeId, Kind: kind}
	rspCh, err := dc.TClient.Send(msg)
	if err != nil {
		return nil, err
	}
	defer close(rspCh)

	rsp := <-rspCh
	if rsp.Method == talk.ERR {
		return nil, errors.New(rsp.Payload.(string))
	} else if rsp.Method == talk.RSP {
		r := rsp.Payload.([]*resource.Resource)
		return r, nil
	} else {
		log.Fatal("No other talk methods should go here!")
		return nil, nil
	}
}

func (dc *DatahubClient) WatchResource(kind resource.Kind, watcher api.ResourceWatcher) error {
	key := talk.TDataKey(dc.NodeId, kind, "#")
	resCh, err := dc.TClient.Watch(key)
	if err != nil {
		return err
	}

	for res := range resCh {
		if res.Sender != dc.TClient.ClientId() { // ignore the change made by client itself
			watcher(res.Payload.(*resource.Resource))
		}
	}
	return nil
}

func (dc *DatahubClient) SaveResource(r *resource.Resource) error {
	msg := talk.NewTMsg()
	msg.Method = talk.PUB
	msg.Key = talk.TDataKey(dc.NodeId, r.Kind, r.Id)
	msg.Payload = r
	_, err := dc.TClient.Send(msg)
	return err
}

func (dc *DatahubClient) DeleteResource(kind resource.Kind, id string) error {
	msg := talk.NewTMsg()
	msg.Method = talk.PUB
	msg.Key = talk.TDataKey(dc.NodeId, kind, id)
	msg.Payload = &resource.Resource{
		ResourceKey: resource.ResourceKey{dc.NodeId, kind, id},
		Ts:          time.Now().UnixNano(),
		Version:     0,
	}
	_, err := dc.TClient.Send(msg)
	return err
}
