package client

import (
	"encoding/binary"
	"fmt"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/datahub/conf"
	"testing"
)

func init() {
	conf.Load("../etc/datahub.yaml")
}

func TestNode(t *testing.T) {
	cli, err := NewDatahubClient(conf.Config.Mqtt, service.BootMan, conf.Config.NodeId)
	if err != nil {
		t.Errorf("fail to create datahub client, err: %s", err.Error())
		return
	}

	res, err := cli.GetResource(resource.KindNode, conf.Config.NodeId)
	if err != nil {
		t.Errorf("fail to get node resource, err: %s", err.Error())
		return
	}

	t.Logf("get node resource: %+v", res)
}

func TestIntByte(t *testing.T) {
	var vi uint64 = 10000000000000
	data := make([]byte,8)
	binary.BigEndian.PutUint64(data, vi)
	vo := binary.BigEndian.Uint64(data)
	fmt.Println(vo)
}