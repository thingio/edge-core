package talk

import (
	"fmt"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/toolkit"
	"strings"
)

type Method string

const (
	ERR Method = "ERR"
	REQ Method = "REQ"
	RSP Method = "RSP"
	PUB Method = "PUB"
	//SUB Method = "SUB"
)

type TMsg struct {
	ID      string
	Sender  string
	Method  Method
	Key     string
	Payload interface{}
}

func NewTMsg() *TMsg {
	return &TMsg{ID: toolkit.NewUUID()}
}

const TChanPrefix = "/TCHAN"
const TDataPrefix = "/TDATA"

func TChanKey(nodeId string, service service.ServiceId, function service.ServiceFunction) string {
	return TChanPrefix + fmt.Sprintf("/%s/%s/%s", nodeId, service, function)
}

func TDataKey(nodeId string, kind resource.Kind, id string) string {
	if kind == resource.KindAny {
		return TDataPrefix + fmt.Sprintf("/%s/%s", nodeId, kind)
	} else {
		return TDataPrefix + fmt.Sprintf("/%s/%s/%s", nodeId, kind, id)
	}
}

type KeyPartIdx int

const (
	KPI_NodeID KeyPartIdx = 2
	// if Method=REQ/RSP/ERR, then Key=/tchan/{node}/{service}/{function}, for a pair of REQ/RSP or REQ/ERR, their Id will be the same
	KPI_Service  KeyPartIdx = 3
	KPI_Function KeyPartIdx = 4
	// if Method=PUB, then Key=/tdata/{node}/{resource_kind}/{resouce_id}
	KPI_ResourceKind KeyPartIdx = 3
	KPI_ResourceID   KeyPartIdx = 4
)

func (t *TMsg) KeyPart(idx KeyPartIdx) string {
	paths := strings.Split(t.Key, "/")
	if len(paths) > int(idx) {
		return paths[idx]
	}
	return ""
}
