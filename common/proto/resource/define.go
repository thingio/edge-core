package resource

type Kind string

const (
	KindAny      Kind = "#"
	KindNode     Kind = "NODE"
	KindPipeline Kind = "PIPELINE"
	KindTask     Kind = "TASK"
	KindModel    Kind = "MODEL"
	KindFunction Kind = "FUNCTION"
	KindDevice   Kind = "DEVICE"
)

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
	Status interface{} `json:"status,omitempty"`
	Ts     int64       `json:"ts,omitempty"`
}
