package resource

type Node struct {
	Id     string `json:"id,omitempty"`
	Os     string `json:"os,omitempty"`
	Kernel string `json:"kernel,omitempty"`
	Arch   string `json:"arch,omitempty"`
}

func (t Node) GetId() string   { return t.Id }
func (t Node) SetId(id string) { t.Id = id }

type Servlet struct {
	Id         string            `json:"id,omitempty"`
	Ns         string            `json:"ns,omitempty"` // namespace
	Image      string            `json:"image,omitempty"`
	Cmd        []string          `json:"cmd,omitempty"`
	Resources  *ServletResources `json:"resources,omitempty"`
	Network    string            `json:"network,omitempty"`
	Privileged bool              `json:"privileged,omitempty"`
	Ports      []string          `json:"ports,omitempty"`
	Devices    []string          `json:"devices,omitempty"`
	Envs       map[string]string `json:"envs,omitempty"`
	Volumes    map[string]string `json:"volumes,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

func (t Servlet) GetId() string   { return t.Id }
func (t Servlet) SetId(id string) { t.Id = id }

type ServletResources struct {
	GPUs     int64   `json:"gpu,omitempty"`      // GPU count, 0 for none
	CPUs     float32 `json:"cpus,omitempty"`     // CPU quota in units of CPUs
	MemLimit int64   `json:"memlimit,omitempty"` // Memory limit in units of Bytes
}

type Funclet struct {
	Id       string `json:"id,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Runtime  string `json:"runtime,omitempty"`
	Code     string `json:"code,omitempty"`
}

func (t Funclet) GetId() string   { return t.Id }
func (t Funclet) SetId(id string) { t.Id = id }

type Applet struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Arch  string `json:"arch,omitempty"`
	Image string `json:"image,omitempty"`
}

func (t Applet) GetId() string   { return t.Id }
func (t Applet) SetId(id string) { t.Id = id }

type Pipeline struct {
	Id    string      `json:"id,omitempty"`
	Genus string      `json:"genus,omitempty"` // mm or ts
	Name  string      `json:"name,omitempty"`
	Mode  string      `json:"mode,omitempty"` // graph or script
	Body  interface{} `json:"body,omitempty"` // pipegraph or string
	Binds []*PipeBind `json:"binds,omitempty"`
}

func (t Pipeline) GetId() string   { return t.Id }
func (t Pipeline) SetId(id string) { t.Id = id }

type PipeBind struct {
	Id    string            `json:"id,omitempty"`    // bind id: if Pipeline.Mode=graph then id=PipeNode.Id
	Name  string            `json:"name,omitempty"`  // bind name
	Type  string            `json:"type,omitempty"`  // bind type： [device|applet|funclet|string]
	Infos map[string]string `json:"infos,omitempty"` // bind info: {product:xxx} when Type=device; {type:xxx} when Type=applet;
	Value string            `json:"value,omitempty"` // bind value: works as a final value in PipeTask, and default value in Pipeline.Binds
	// TIPS for web-dev:
	// case Type=device then GET /data/devices?product=xxx for user to select
	// case Type=applet then GET /data/applets?type=xxx to for user to select
	// case Type=funclet then GET /data/funclets for user to select
	// case Type=string then user can input anything
}

type PipeGraph struct {
	Nodes []*PipeNode `json:"nodes,omitempty"`
	Links []*PipeLink `json:"links,omitempty"`
	Ui    string      `json:"ui,omitempty"`
}

// PipeWidget is the template for PipeNode in PipeGraph
type PipeWidget struct {
	Id         string   `json:"id,omitempty"`
	Genus      string   `json:"genus,omitempty"`
	Group      string   `json:"group,omitempty"`
	Name       string   `json:"name,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	BindType   string   `json:"bind_type,omitempty"` // PipeWidget.BindType -> PipeNode.BindType -> Pipeline.Args.Type
	Params     []*Param `json:"params,omitempty"`    // PipeWidget.Params will be used to generates Pipeline.Nodes[i].Props
	Upstream   []string `json:"upstream,omitempty"`
	Downstream []string `json:"downstream,omitempty"`
}

func (t PipeWidget) GetId() string   { return t.Id }
func (t PipeWidget) SetId(id string) { t.Id = id }

// PipeNode represents a node in PipeGraph
type PipeNode struct {
	NodeId    string            `json:"node_id,omitempty"` // user generated uuid
	NodeName  string            `json:"node_name,omitempty"`
	WidgetId  string            `json:"widget_id,omitempty"`
	BindType  string            `json:"bind_type,omitempty"` // PipeWidget.BindType -> PipeNode.BindType -> Pipeline.Binds[i].Type
	NodeProps map[string]string `json:"node_props,omitempty"`
	NodeUi    string            `json:"node_ui,omitempty"` // such as x,y position in a 2-D graph
}

// PipeLink represents link between 2 PipeNodes
type PipeLink struct {
	LinkId     string `json:"link_id,omitempty"` // user generated uuid
	FromNodeId string `json:"from_node_id,omitempty"`
	ToNodeId   string `json:"to_node_id,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	LinkUi     string `json:"link_ui,omitempty"`
}

type PipeTask struct {
	Id         string      `json:"id,omitempty"`          // task id
	PipelineId string      `json:"pipeline_id,omitempty"` // related pipeline id
	Genus      string      `json:"genus,omitempty"`       // related pipeline type: mm or ts
	Form       string      `json:"form,omitempty"`        // related pipeline form: graph or script
	Name       string      `json:"name,omitempty"`
	Body       interface{} `json:"body,omitempty"`
	Binds      []*PipeBind `json:"binds,omitempty"`
}

func (t PipeTask) GetId() string   { return t.Id }
func (t PipeTask) SetId(id string) { t.Id = id }

type Device struct {
	Id        string            `json:"id,omitempty"`
	Genus     string            `json:"genus,omitempty"` // mm or ts
	Name      string            `json:"name,omitempty"`
	ProductId string            `json:"product_id,omitempty"`
	Props     map[string]string `json:"props,omitempty"`
}

func (t Device) GetId() string   { return t.Id }
func (t Device) SetId(id string) { t.Id = id }

type DeviceProduct struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Genus      string `json:"genus,omitempty"`
	Desc       string `json:"desc,omitempty"`
	ProtocolId string `json:"protocol_id,omitempty"`
}

func (t DeviceProduct) GetId() string   { return t.Id }
func (t DeviceProduct) SetId(id string) { t.Id = id }

type DeviceProtocol struct {
	Id     string   `json:"id,omitempty"`
	Name   string   `json:"name,omitempty"`
	Genus  string   `json:"genus,omitempty"`
	Desc   string   `json:"desc,omitempty"`
	Params []*Param `json:"params,omitempty"` // DeviceProtocol.Params will be use to generate Device.Props
}

func (t DeviceProtocol) GetId() string   { return t.Id }
func (t DeviceProtocol) SetId(id string) { t.Id = id }

type Param struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Type     string `json:"type,omitempty"`
	Style    string `json:"style,omitempty"`
	Default  string `json:"default,omitempty"`
	Range    string `json:"range,omitempty"` // eg. "true@@自动模式,false@@手动模式"
	Required bool   `json:"required,omitempty"`
}

type Alert struct {
	Id      string `json:"id,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
	Image   string `json:"image,omitempty"` // better be []byte, but after json serialization, they will be base64 anyway
}

func (t Alert) GetId() string   { return t.Id }
func (t Alert) SetId(id string) { t.Id = id }
