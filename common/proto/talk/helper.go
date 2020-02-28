package talk

import (
	"fmt"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/service"
	"github.com/thingio/edge-core/common/toolkit"
	"strings"
)

func NewTMessage() *TMessage {
	return &TMessage{Id: toolkit.NewUUID()}
}

const TMessageChatPrefix = "/TCHAT"
const TMessageDataPrefix = "/TDATA"

// if Method=MethodREQ/MethodRSP/MethodERR, then Key=/tchan/{node}/{service}/{function}, for a pair of MethodREQ/MethodRSP or MethodREQ/MethodERR, their GetId will be the same
func TMessageChatKey(nodeId string, service service.ServiceId, function service.ServiceFunction) string {
	return TMessageChatPrefix + fmt.Sprintf("/%s/%s/%s", nodeId, service, function)
}

// if Method=MethodPUB, then Key=/tdata/{node}/{resource_kind}/{resouce_id}
func TMessageDataKey(nodeId string, kind *resource.Kind, id string) string {
	if kind == resource.KindAny {
		return TMessageDataPrefix + fmt.Sprintf("/%s/%s", nodeId, kind.Name)
	} else {
		return TMessageDataPrefix + fmt.Sprintf("/%s/%s/%s", nodeId, kind.Name, id)
	}
}



const (
	MethodERR = "ERR"
	MethodREQ = "REQ"
	MethodRSP = "RSP"
	MethodPUB = "PUB"
	//MethodSUB Method = "SUB"
)

const (
	keyPartNodeId       = 2
	keyPartService      = 3
	keyPartFunction     = 4
	keyPartResourceKind = 3
	keyPartResourceId   = 4
)
func (t *TMessage) NodeId() string {
	return t.keyPart(keyPartNodeId)
}
func (t *TMessage) ServiceId() string {
	return t.keyPart(keyPartService)
}
func (t *TMessage) ServiceFunction() string {
	return t.keyPart(keyPartFunction)
}
func (t *TMessage) ResourceKind() string {
	return t.keyPart(keyPartResourceKind)
}
func (t *TMessage) ResourceId() string {
	return t.keyPart(keyPartResourceId)
}
func (t *TMessage) keyPart(idx int) string {
	paths := strings.Split(t.Key, "/")
	if len(paths) > idx {
		return paths[idx]
	}
	return ""
}
