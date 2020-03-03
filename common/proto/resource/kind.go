package resource

import (
	"github.com/thingio/edge-core/common/toolkit"
)

type Kind struct {
	Name         string
	Stateful     bool
	Cloneable     bool
	SampleObject interface{}
	NewObject    func(id string) IdObject
}

var (
	KindAny      = &Kind{Name: "#", Cloneable: false, Stateful: false, SampleObject: nil, NewObject: func(id string) IdObject { return nil }}
	KindNode     = &Kind{Name: "node", Cloneable: false, Stateful: true, SampleObject: Node{}, NewObject: func(id string) IdObject { return &Node{Id: id} }}
	KindDevice   = &Kind{Name: "device", Cloneable: true, Stateful: false, SampleObject: Device{}, NewObject: func(id string) IdObject { return &Device{Id: id, Props: make(map[string]string, 0)} }}
	KindPipeline = &Kind{Name: "pipeline", Cloneable: true, Stateful: false, SampleObject: Pipeline{}, NewObject: func(id string) IdObject { return &Pipeline{Id: id, Body: &PipeGraph{}, Binds: make(map[string]*PipeBind, 0)} }}
	KindPipeTask = &Kind{Name: "pipetask", Cloneable: true, Stateful: true, SampleObject: PipeTask{}, NewObject: func(id string) IdObject { return &PipeTask{Id: id, Binds: make(map[string]*PipeBind, 0)} }}
	KindApplet   = &Kind{Name: "applet", Cloneable: true, Stateful: true, SampleObject: Applet{}, NewObject: func(id string) IdObject { return &Applet{Id: id} }}
	KindFunclet  = &Kind{Name: "funclet", Cloneable: true, Stateful: true, SampleObject: Funclet{}, NewObject: func(id string) IdObject { return &Funclet{Id: id} }}
	KindServlet  = &Kind{Name: "servlet", Cloneable: false, Stateful: true, SampleObject: Servlet{}, NewObject: func(id string) IdObject { return &Servlet{Id: id, Envs: make(map[string]string, 0), Volumes: make(map[string]string, 0), Labels: make(map[string]string, 0)}	}}
	KindState    = &Kind{Name: "state", Cloneable: false, Stateful: false, SampleObject: State{}, NewObject: func(id string) IdObject { return State{"id": id} }}
	KindProduct  = &Kind{Name: "product", Cloneable: false, Stateful: false, SampleObject: DeviceProduct{}, NewObject: func(id string) IdObject { return &DeviceProduct{Id: id} }}
	KindProtocol = &Kind{Name: "protocol", Cloneable: false, Stateful: false, SampleObject: DeviceProtocol{}, NewObject: func(id string) IdObject { return &DeviceProtocol{Id: id, Params: make([]*DeviceParam, 0)} }}
	KindWidget   = &Kind{Name: "widget", Cloneable: false, Stateful: false, SampleObject: PipeWidget{}, NewObject: func(id string) IdObject { return &PipeWidget{Id: id, Params: make([]*PipeParam, 0)} }}
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

