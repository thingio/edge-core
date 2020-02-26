package resource

type Node struct {
	Id        string            `json:"id,omitempty"`
	Os        string            `json:"os,omitempty"`
	Kernel    string            `json:"kernel,omitempty"`
	Arch      string            `json:"arch,omitempty"`
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
	Id      string    `json:"id,omitempty"`
	Genus   string    `json:"genus,omitempty"` // mm or ts
	Name    string    `json:"name,omitempty"`
	BodyDef string    `json:"body_def,omitempty"`
	ArgDefs []ArgBind `json:"arg_defs,omitempty"`
}

func (t Pipeline) GetId() string   { return t.Id }
func (t Pipeline) SetId(id string) { t.Id = id }

type PipeTask struct {
	Id         string            `json:"id,omitempty"`          // task id
	PipelineId string            `json:"pipeline_id,omitempty"` // related pipeline id
	Genus      string            `json:"genus,omitempty"`       // related pipeline type: mm or ts
	Name       string            `json:"name,omitempty"`
	Body       string            `json:"body,omitempty"`
	Args       map[string]string `json:"args,omitempty"`
}

func (t PipeTask) GetId() string   { return t.Id }
func (t PipeTask) SetId(id string) { t.Id = id }

type ArgBind struct {
	Id   string `json:"id,omitempty"`   // arg id
	Name string `json:"name,omitempty"` // arg name
	Type string `json:"type,omitempty"` // arg type for frontend to render corresponding input style
}

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
	Id         string       `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Genus      string       `json:"genus,omitempty"`
	Desc       string       `json:"desc,omitempty"`
	PropFields []*PropField `json:"prop_fields,omitempty"` // DeviceProtocol.PropFields will be use to generate Device.Props
}

func (t DeviceProtocol) GetId() string   { return t.Id }
func (t DeviceProtocol) SetId(id string) { t.Id = id }

type PropField struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Type     string `json:"type,omitempty"`
	Style    string `json:"style,omitempty"`
	Default  string `json:"default,omitempty"`
	Range    string `json:"range,omitempty"` //"true@@自动模式,false@@手动模式"
	Required bool   `json:"required,omitempty"`
}
