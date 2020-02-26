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
	"github.com/thingio/edge-core/datahub/storage"
	"time"
)

type DatahubServer struct {
	*talk.TClient
	rdb    storage.ResourceStorage
	nodeId string
}

func NewDatahubServer(config conf.DatahubConfig) *DatahubServer {
	s := new(DatahubServer)
	s.TClient = talk.NewTClient(config.Mqtt, string(service.DataHub), talkpb.TMessageChatPrefix+"/"+config.NodeId)
	s.nodeId = config.NodeId
	s.rdb = storage.NewResourceStorage(config.DB)
	s.SetReqHandler(s.Handle)
	return s
}

var (
	ResourceDirs = map[*resource.Kind]string{
		resource.KindProduct:  "etc/products/",
		resource.KindProtocol: "etc/protocols/",
	}
)

func (t *DatahubServer) Init() error {
	// load products and protocols resource into rdb on init
	for kind, path := range ResourceDirs {
		log.Infof("loading '%s' resources from '%s'", kind, path)
		rs, err := resource.LoadResourcesFromDir(kind, path)
		if err != nil {
			return err
		}
		log.Infof("%d '%s' resources loaded from '%s'", len(rs), kind, path)
		for _, r := range rs {
			r.NodeId = t.nodeId
			if err = t.rdb.Put(r); err != nil {
				return err
			}
		}
		log.Infof("%d '%s' resources saved into '%s'", len(rs), kind, path)
	}

	// load node resource into rdb on init and periodically update node state resource
	r := resource.KindNode.NewResourceOf(NodeData(t.nodeId))
	r.NodeId = t.nodeId
	if err := t.rdb.Put(r); err != nil {
		return err
	}

	go func(t *DatahubServer) {
		ticker := time.NewTicker(10 * time.Second)
		for {
			<- ticker.C
			r := resource.KindState.NewResourceOf(NodeState(t.nodeId))
			r.NodeId = t.nodeId
			if err := t.rdb.Put(r); err != nil {
				log.WithError(err).Warnf("fail to save node state to rdb")
			}
		}
	}(t)

	return nil
}


func (t *DatahubServer) Start() error {
	if err := t.rdb.Init(); err != nil {
		return err
	}

	if err := t.Init(); err != nil {
		return err
	}
	if err := t.Connect(); err != nil {
		return err
	}

	return t.Serve()
}

// LOOP: keep handling incoming resource events
func (t *DatahubServer) Serve() error {

	tchanKey := talkpb.TMessageDataKey(conf.Config.NodeId, resource.KindAny, "")
	if ch, err := t.Watch(tchanKey); err != nil {
		return err
	} else {
		for tmsg := range ch {
			kind := resource.KindOf(tmsg.ResourceKind())
			if err := t.Save(kind, tmsg.Payload); err != nil {
				log.Infof("fail to save %s resource: %s", kind, tmsg.Payload)
			}
		}
	}
	return nil
}

func (t *DatahubServer) Handle(function string, payload []byte) ([]byte, error) {

	req := resource.Key{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.Wrap(err, errors.Errorf("REQ payload can only be unmarshalled as resource.Key, payload: %s", payload))
	}

	switch service.ServiceFunction(function) {
	case service.FuncGet:
		rsp, err := t.GetResource(req)
		if err != nil {
			return nil, err
		}
		return resource.MarshalResource(rsp)
	case service.FuncList:
		rsp, err := t.ListResource(req)
		if err != nil {
			return nil, err
		}
		return resource.MarshalResourceList(rsp)
	case service.FuncState:
		//TODO
		return nil, errors.NotImplementedf("function '%s'", function)
	default:
		return nil, errors.NotSupportedf("function '%s'", function)
	}

}

func (t *DatahubServer) Save(kind *resource.Kind, data []byte) error {
	log.Infof("receive '%s' resource event: %s", kind, data)
	r, err := resource.UnmarshalResource(kind, data)
	if err != nil {
		return err
	}
	return t.rdb.Put(r)
}

func (t *DatahubServer) GetResource(key resource.Key) (*resource.Resource, error) {
	log.Infof("receive get call: %+v", key)
	return t.rdb.Get(key)
}

func (t *DatahubServer) ListResource(key resource.Key) ([]*resource.Resource, error) {
	log.Infof("receive list call: %+v", key)
	return t.rdb.List(key)
}
