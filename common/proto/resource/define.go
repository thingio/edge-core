package resource

import (
	"github.com/thingio/edge-core/common/toolkit"
	"log"
)

type Kind string

const (
	KindAny      Kind = "#"
	KindNode     Kind = "node"
	KindDevice   Kind = "device"
	KindPipeline Kind = "pipeline"
	KindPipeTask Kind = "pipetask"
	KindApplet   Kind = "applet"
	KindFunclet  Kind = "funclet"
	KindServlet  Kind = "servlet"
)

var AllKinds = []Kind{KindNode, KindPipeline, KindPipeTask, KindApplet, KindFunclet, KindDevice, KindServlet}

func (k Kind) NewId() string {
	return string(k) + "-" + toolkit.NewUUID()
}

type IdObject interface {
	GetId() string
	SetId(id string)
}

func (k Kind) NewObject(id string) IdObject {
	switch k {
	case KindNode:
		return Node{Id: id}
	case KindPipeline:
		return Pipeline{Id: id}
	case KindPipeTask:
		return PipeTask{Id: id}
	case KindApplet:
		return Applet{Id: id}
	case KindFunclet:
		return Funclet{Id: id}
	case KindServlet:
		return Servlet{Id: id}
	case KindDevice:
		return Device{Id: id}
	default:
		log.Fatalf("%s not support yet\n", k)
	}
	return nil
}

func (k Kind) NewEmptyResource() *Resource {
	return &Resource{Value: k.NewObject("")}
}

func (k Kind) NewResource() *Resource {
	id := k.NewId()
	return k.NewResourceWithId(id)
}

func (k Kind) NewResourceWithId(id string) *Resource {
	obj := k.NewObject(id)
	key := ResourceKey{Kind: k, Id: id}
	return &Resource{ResourceKey: key, Value: obj, Ts: toolkit.Now(), Version: 1}
}

type ResourceKey struct {
	NodeId string `json:"node_id,omitempty"`
	Kind   Kind   `json:"kind,omitempty"`
	Id     string `json:"id,omitempty"`
}

type Resource struct {
	ResourceKey
	Value   interface{} `json:"value,omitempty"`
	Ts      int64       `json:"ts,omitempty"`
	Version int64       `json:"version,omitempty"`
}

type ResourceStatus struct {
	ResourceKey
	Status IdObject `json:"status,omitempty"`
	Ts     int64    `json:"ts,omitempty"`
}
