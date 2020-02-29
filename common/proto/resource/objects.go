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
	Mode  string `json:"mode,omitempty"` // unmanaged (remote service) or managed (started by bootman)
	Url   string `json:"url,omitempty"`  // http://host:port for unmanaged applets
	Specs struct {
		Path   string `json:"path,omitempty"`
		Method string `json:"method,omitempty"`
		Req    string `json:"req,omitempty"`
	} `json:"specs,omitempty"`
	Configs interface{} `json:"configs,omitempty"` // for managed applets, user will provide configs for bootman
}

func (t Applet) GetId() string   { return t.Id }
func (t Applet) SetId(id string) { t.Id = id }

type Pipeline struct {
	Id    string               `json:"id,omitempty"`
	Genus string               `json:"genus,omitempty"` // mm or ts
	Name  string               `json:"name,omitempty"`
	Mode  string               `json:"mode,omitempty"` // graph or script
	Body  *PipeGraph           `json:"body,omitempty"` // TODO: pipegraph or string
	Binds map[string]*PipeBind `json:"binds,omitempty"`
}

func (t Pipeline) GetId() string   { return t.Id }
func (t Pipeline) SetId(id string) { t.Id = id }

type PipeBind struct {
	Id           string `json:"id,omitempty"` // bind id: if Pipeline.Mode=graph then id={PipeNode.Id}.{Param.Id}
	NodeId       string `json:"node_id,omitempty"`
	ParamId      string `json:"param_id,omitempty"`
	NodeName     string `json:"node_name,omitempty"`
	ParamName    string `json:"param_name,omitempty"`
	Value        string `json:"value,omitempty"`         // bind value: final value in PipeTask
	DefaultValue string `json:"default_value,omitempty"` // bind default value: default value in Pipeline.Binds

	Style ParamStyle `json:"type,omitempty"`  // style of this bind: [roi|input|selector|dynamic_selector]
	Range ParamRange `json:"range,omitempty"` // value range of this bind
}

type PipeGraph struct {
	Nodes []*PipeNode `json:"nodes,omitempty"`
	Links []*PipeLink `json:"links,omitempty"`
	Ui    interface{} `json:"ui,omitempty"` // any UI data, won't be used by backend
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
	NodeProps map[string]string `json:"node_props,omitempty"`
	NodeUi    interface{}       `json:"node_ui,omitempty"` // any UI data, such as x,y position in a 2D plate
}

// PipeLink represents link between 2 PipeNodes
type PipeLink struct {
	LinkId     string      `json:"link_id,omitempty"` // user generated uuid
	FromNodeId string      `json:"from_node_id,omitempty"`
	ToNodeId   string      `json:"to_node_id,omitempty"`
	LinkUi     interface{} `json:"link_ui,omitempty"` // any UI data, won't be used by backend
}

type PipeTask struct {
	Id         string               `json:"id,omitempty"`          // task id
	PipelineId string               `json:"pipeline_id,omitempty"` // related pipeline id
	Genus      string               `json:"genus,omitempty"`       // related pipeline type: mm or ts
	Enable     bool                 `json:"enable,omitempty"`
	Name       string               `json:"name,omitempty"`
	Binds      map[string]*PipeBind `json:"binds,omitempty"`
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

type ParamScope = string

type ParamStyle = string

// ParamRange specifies the range of values of different ParamStyle:
// case Style=input: user can input anything
// case Style=selector: user can select an item like "true" in bind range "true@@是,false@@否"
// case Style=dynamic_selector: user can select an item from response of api specified by bind range like "id,name@/api/v1/res/funclets"
type ParamRange = string

const (
	All      ParamScope = "*"
	PipeOnly ParamScope = "pipe"
	TaskOnly ParamScope = "task"

	ROI             ParamStyle = "roi"
	Input           ParamStyle = "input"
	Textarea        ParamStyle = "textarea"
	Selector        ParamStyle = "selector"
	DynamicSelector ParamStyle = "dynamic_selector"
)

type Param struct {
	Id       string     `json:"id,omitempty"`
	Name     string     `json:"name,omitempty"`
	Desc     string     `json:"desc,omitempty"`
	Type     string     `json:"type,omitempty"`
	Default  string     `json:"default,omitempty"`
	Required bool       `json:"required,omitempty"`
	Depends  []string   `json:"depends,omitempty"` // param ids that this parameter depends on
	Scope    ParamScope `json:"scope,omitempty"`
	Style    ParamStyle `json:"style,omitempty"`
	Range    ParamRange `json:"range,omitempty"`
}

type Alert struct {
	Id         string `json:"id,omitempty"`
	PipeTaskId string `json:"pipetask_id,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Message    string `json:"message,omitempty"`
	Level      string `json:"level,omitempty"`
	Data       string `json:"data,omitempty"`  // better be []byte, but after json serialization, they will be base64 anyway
	Image      string `json:"image,omitempty"` // same as above
	Ts         string `json:"ts,omitempty"`
}
