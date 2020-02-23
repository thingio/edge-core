package resource

import (
	"fmt"
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

func (k Kind) NewSample(id string) IdObject {
	switch k {
	case KindNode:
		return Node{Id: id}
	case KindPipeline:
		return Pipeline{Id: id, ArgDefs: make([]ArgBind, 0)}
	case KindPipeTask:
		return PipeTask{Id: id, Args: make(map[string]string, 0)}
	case KindApplet:
		return Applet{Id: id}
	case KindFunclet:
		return Funclet{Id: id}
	case KindServlet:
		return Servlet{Id: id, Envs: make(map[string]string, 0), Volumes: make(map[string]string, 0), Labels: make(map[string]string, 0)}
	case KindDevice:
		return Device{Id: id, Props: make(map[string]string, 0)}
	default:
		log.Fatalf("%s not support yet\n", k)
	}
	return nil
}

func (k Kind) NewObject(id string) IdObject {
	switch k {
	case KindNode:
		return &Node{Id: id}
	case KindPipeline:
		return &Pipeline{Id: id, ArgDefs: make([]ArgBind, 0)}
	case KindPipeTask:
		return &PipeTask{Id: id, Args: make(map[string]string, 0)}
	case KindApplet:
		return &Applet{Id: id}
	case KindFunclet:
		return &Funclet{Id: id}
	case KindServlet:
		return &Servlet{Id: id, Envs: make(map[string]string, 0), Volumes: make(map[string]string, 0), Labels: make(map[string]string, 0)}
	case KindDevice:
		return &Device{Id: id, Props: make(map[string]string, 0)}
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
	key := Key{Kind: k, Id: id}
	return &Resource{Key: key, Value: obj, Ts: toolkit.Now(), Version: 1}
}

type Key struct {
	NodeId string `json:"node_id,omitempty"`
	Kind   Kind   `json:"kind,omitempty"`
	Id     string `json:"id,omitempty"`
}

func (k Key) String() string {
	return fmt.Sprintf("/%s/%s/%s", k.NodeId, k.Kind, k.Id)
}

type Resource struct {
	Key
	Value   interface{} `json:"value,omitempty"`
	Ts      int64       `json:"ts,omitempty"`
	Version int64       `json:"version,omitempty"`
}

type ResourceStatus struct {
	Key
	Status IdObject `json:"status,omitempty"`
	Ts     int64    `json:"ts,omitempty"`
}
