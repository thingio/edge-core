package resource

type Node struct {
	Id        string            `json:"id,omitempty"`
	Os        string            `json:"os,omitempty"`
	Kernel    string            `json:"kernel,omitempty"`
	Arch      string            `json:"arch,omitempty"`
	BootTime  int64             `json:"boot_time,omitempty"`
	LocalTime string            `json:"local_time,omitempty"`
	Stats     map[string]string `json:"stats,omitempty"`
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

type Device struct {
	Id      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Product string            `json:"product,omitempty"`
	Props   map[string]string `json:"Props,omitempty"`
}

func (t Device) GetId() string   { return t.Id }
func (t Device) SetId(id string) { t.Id = id }


type Pipeline struct {
	Id       string    `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	Name     string    `json:"name,omitempty"`
	BodyDef  string    `json:"body_def,omitempty"`
	ArgDefs  []ArgBind `json:"arg_defs,omitempty"`
}

func (t Pipeline) GetId() string   { return t.Id }
func (t Pipeline) SetId(id string) { t.Id = id }

type PipeTask struct {
	Id           string            `json:"id,omitempty"`
	PipelineId   string            `json:"pipeline_id,omitempty"`
	PipelineType string            `json:"type,omitempty"`
	Name         string            `json:"name,omitempty"`
	Body         string            `json:"body,omitempty"`
	Args         map[string]string `json:"args,omitempty"`
}

func (t PipeTask) GetId() string   { return t.Id }
func (t PipeTask) SetId(id string) { t.Id = id }

type ArgBind struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
