package server

import (
	"github.com/juju/errors"
	"github.com/thingio/edge-core/datahub/conf"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/talk"
	"github.com/thingio/edge-core/common/toolkit"
	"github.com/thingio/edge-core/datahub/resource/api"
)

func Start() error {
	cli := talk.NewTClient(conf.Config.Mqtt, service.DeviceMan)
	if err := cli.Connect(); err != nil {
		return err
	}
	cli.SetReqHandler(reqHandler)

	if ch, err := cli.Watch(talk.DataTopicPrefix + "#"); err != nil {
		return err
	} else {
		for tmsg := range ch {
			kind := api.Kind(tmsg.KeyPart(talk.IdxResourceKind))
			Save(kind, tmsg.Payload)
		}
	}

	return nil
}

func reqHandler(path string, payload interface{}) (i interface{}, e error) {
	if path == "node" {
		return toolkit.NodeData(conf.Config.NodeId), nil
	}
	return nil, errors.NotImplementedf("datahub kind: device/pipeline/task/model/function")
}

func Save(kind api.Kind, data interface{}) {

}
