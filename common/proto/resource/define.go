package resource

import (
	"fmt"
	"github.com/thingio/edge-core/common/toolkit"
)

const (
	GenusTS = "ts" // 时序分类
	GenusMM = "mm" // 多媒体分类
)

type IdObject interface {
	GetId() string
	SetId(id string)
}

type Kind struct {
	Name         string
	Stateful     bool
	SampleObject interface{}
	NewObject    func(id string) IdObject
}

var (
	KindAny      = &Kind{Name: "#", Stateful: false, SampleObject: nil, NewObject: func(id string) IdObject { return nil }}
	KindNode     = &Kind{Name: "node", Stateful: true, SampleObject: Node{}, NewObject: func(id string) IdObject { return &Node{Id: id} }}
	KindDevice   = &Kind{Name: "device", Stateful: false, SampleObject: Device{}, NewObject: func(id string) IdObject { return &Device{Id: id, Props: make(map[string]string, 0)} }}
	KindPipeline = &Kind{Name: "pipeline", Stateful: false, SampleObject: Pipeline{}, NewObject: func(id string) IdObject { return &Pipeline{Id: id, Body: &PipeGraph{}, Binds: make(map[string]*PipeBind, 0)} }}
	KindPipeTask = &Kind{Name: "pipetask", Stateful: true, SampleObject: PipeTask{}, NewObject: func(id string) IdObject { return &PipeTask{Id: id, Binds: make(map[string]*PipeBind, 0)} }}
	KindApplet   = &Kind{Name: "applet", Stateful: true, SampleObject: Applet{}, NewObject: func(id string) IdObject { return &Applet{Id: id} }}
	KindFunclet  = &Kind{Name: "funclet", Stateful: true, SampleObject: Funclet{}, NewObject: func(id string) IdObject { return &Funclet{Id: id} }}
	KindServlet  = &Kind{Name: "servlet", Stateful: true, SampleObject: Servlet{}, NewObject: func(id string) IdObject { return &Servlet{Id: id, Envs: make(map[string]string, 0), Volumes: make(map[string]string, 0), Labels: make(map[string]string, 0)}	}}
	KindState    = &Kind{Name: "state", Stateful: false, SampleObject: State{}, NewObject: func(id string) IdObject { return State{"id": id} }}
	KindProduct  = &Kind{Name: "product", Stateful: false, SampleObject: DeviceProduct{}, NewObject: func(id string) IdObject { return &DeviceProduct{Id: id} }}
	KindProtocol = &Kind{Name: "protocol", Stateful: false, SampleObject: DeviceProtocol{}, NewObject: func(id string) IdObject { return &DeviceProtocol{Id: id, Params: make([]*Param, 0)} }}
	KindWidget   = &Kind{Name: "widget", Stateful: false, SampleObject: PipeWidget{}, NewObject: func(id string) IdObject { return &PipeWidget{Id: id, Params: make([]*Param, 0)} }}
)

var AllKinds = []*Kind{
	KindNode, KindPipeline, KindPipeTask,
	KindDevice, KindProduct, KindProtocol,
	KindApplet, KindFunclet, KindServlet,
	KindState, KindWidget,
}

var KindMap = make(map[string]*Kind)

func init() {
	for _, k := range AllKinds {
		KindMap[k.Name] = k
	}
}

func KindOf(name string) *Kind {
	return KindMap[name]
}

func (k *Kind) NewId() string {
	return k.Name + "-" + toolkit.NewUUID()
}

func (k *Kind) String() string {
	return k.Name
}

func (k *Kind) NewEmptyResource() *Resource {
	return &Resource{Value: k.NewObject("")}
}

func (k *Kind) NewResource() *Resource {
	id := k.NewId()
	return k.NewResourceWithId(id)
}

func (k *Kind) NewResourceWithId(id string) *Resource {
	obj := k.NewObject(id)
	key := Key{Kind: k.Name, Id: id}
	return &Resource{Key: key, Value: obj, Ts: toolkit.Now(), Version: 1}
}

func (k *Kind) NewResourceOf(obj IdObject) *Resource {
	key := Key{Kind: k.Name, Id: obj.GetId()}
	return &Resource{Key: key, Value: obj, Ts: toolkit.Now(), Version: 1}
}

type Key struct {
	NodeId string `json:"node_id,omitempty"`
	Kind   string `json:"kind,omitempty"`
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

type State map[string]string

func (t State) GetId() string   { return t["id"] }
func (t State) SetId(id string) { t["id"] = id }

type ResourceState struct {
	Key
	Value   interface{} `json:"value,omitempty"`
	Ts      int64       `json:"ts,omitempty"`
	State   *State      `json:"state,omitempty"`
	StateTs int64       `json:"state_ts,omitempty"`
}

func MakeResourceState(value *Resource, state *Resource) *ResourceState {
	if state == nil {
		return &ResourceState{
			Key:   value.Key,
			Value: value.Value,
			Ts:    value.Ts,
		}
	} else {
		return &ResourceState{
			Key:     value.Key,
			Value:   value.Value,
			Ts:      value.Ts,
			State:   state.Value.(*State),
			StateTs: state.Ts,
		}
	}

}
