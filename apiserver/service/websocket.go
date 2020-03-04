package service

import (
	"github.com/googollee/go-socket.io"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/data"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/datahub/api"
	"net/http"
	"strings"
)

const (
	SocketIoNS   = "/"
	SocketIoRoom = "default"

	SocketEventNotif = "notif"
	SocketEventState = "state"
	SocketEventAlert = "alert"
)

var WsApi = &wsApi{}

type wsApi struct {
	cli    api.DatahubApi
	server *socketio.Server
}

func (t *wsApi) Init(cli api.DatahubApi) error {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return err
	}
	t.server = server
	t.cli = cli

	t.WatchServerChange()

	if err := t.WatchState(); err != nil {
		return err
	}
	if err := t.WatchAlert(); err != nil {
		return err
	}

	return nil
}

func (t *wsApi) Handler() http.Handler {
	return t.server
}

func (t *wsApi) WatchServerChange() {
	t.server.OnConnect(SocketIoNS, func(s socketio.Conn) error {
		log.Infof("%s [%s] connected", s.RemoteAddr().String(), s.ID())
		return nil
	})

	t.server.OnError(SocketIoNS, func(s socketio.Conn, e error) {
		log.WithError(e).Infof("%s [%s] meet error", s.RemoteAddr().String(), s.ID())
	})

	t.server.OnDisconnect(SocketIoNS, func(s socketio.Conn, reason string) {
		log.Infof("%s [%s] closed, reason: %s", s.RemoteAddr().String(), s.ID(), reason)
	})
}

func (t *wsApi) WatchState() error {
	return t.cli.WatchResource(resource.KindState, func(s *resource.Resource) {
		kindStr := strings.SplitN(s.Id, "-", 2)[0]
		kind := resource.KindOf(kindStr)

		var result *resource.ResourceState

		r, err := t.cli.GetResource(kind, s.Id)
		if err != nil {
			log.WithError(err).Warnf("watch state: state %s without resource", s.Id)
			result = &resource.ResourceState{
				Key:     s.Key,
				State:   s.Value.(*resource.State),
				StateTs: s.Ts,
			}
			result.Kind = kind.Name
		} else {
			result = resource.MakeResourceState(r, s)
		}

		t.server.BroadcastToRoom(SocketIoNS, SocketIoRoom, SocketEventState, result)
	}, true)
}

func (t *wsApi) WatchAlert() error {
	// TODO
	return nil
}

func (t *wsApi) SendNotif(notif data.Notif) {
	success := t.server.BroadcastToRoom(SocketIoNS, SocketIoRoom, SocketEventState, notif)
	if !success {
		log.Error("fail to send notif, please check")
	}
}
