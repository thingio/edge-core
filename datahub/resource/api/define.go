package api

import "context"

type Kind string

const (
	Node     Kind = "NODE"
	Pipeline Kind = "PIPELINE"
	Task     Kind = "TASK"
	Model    Kind = "MODEL"
	Function Kind = "FUNCTION"
	Device   Kind = "DEVICE"
)

type Action int32

const (
	None         Action = 0
	Create       Action = 1
	Update       Action = 2
	Delete       Action = 3
	Purge        Action = 4
	Connected    Action = -1
	Disconnected Action = -2
)

type Resource struct {
	Kind    Kind        `json:"kind,omitempty"`
	Owner   string      `json:"owner,omitempty"`
	ID      string      `json:"id,omitempty"`
	Spec    interface{} `json:"spec,omitempty"`
	Status  interface{} `json:"status,omitempty"`
	Ts      int64       `json:"ts,omitempty"`
	Version int64       `json:"version,omitempty"`
	Editor    string    `json:"editor,omitempty"`
}

type ResourceApi interface {
	/* GET api: non-deleted resources */
	GetResource(nodeId string, kind Kind, id string) (*Resource, error)
	ListResources(nodeId string, kind Kind) ([]*Resource, error)
	QueryResources(nodeId string, kind Kind, idPrefix string) ([]*Resource, error)

	/* WATCH api: include the ones editor itself made */
	WatchResource(nodeId string, kind Kind) *ResourceWatcher

	/* UPDATE api: save the entire resource */
	SaveResource(r *Resource) (*Resource, error)

	/* DELETE api: soft deletion, which marks version of the resource to 0 */
	DeleteResource(nodeId string, kind Kind, id string) (*Resource, error)
	DeleteResources(nodeId string, kind Kind, idPrefix string) ([]*Resource, error)

	/* PURGE api: hard-delete, remove the resource from database */
	PurgeResource(nodeId string, kind Kind, id string) (*Resource, error)
	PurgeResources(nodeId string, kind Kind, idPrefix string) ([]*Resource, error)

	/* HELPER api: utility functions */
	ListOwners() ([]string, error)
	ListResourceKeys(nodeId string, kind Kind) ([]string, error)
	ExistResource(nodeId string, kind Kind, id string) (bool, error)

	/* SYNC api: internally used when syncing data between edge nodes and cloud nodes */
	GetResourceSafe(nodeId string, kind Kind, id string) (*Resource, error)
	GetResourceTs(nodeId string, kind Kind, id string) (*ResourceTs, error)
	ListResourceTs(nodeId string, kind Kind) ([]*ResourceTs, error)

}

type ResourceWatcher struct {
	Key        string
	EventQueue chan *ResourceEvent
	Context    context.Context
	Cancel     context.CancelFunc
}

type ResourceEvent struct {
	Key          string    `json:"key,omitempty"`
	Action    	 Action    `json:"action,omitempty"`
	Editor       string    `json:"editor,omitempty"`
	Resource     *Resource `json:"resource,omitempty"`
	PrevResource *Resource `json:"prev_resource,omitempty"`
	Error        error     `json:"error,omitempty"`
}

type ResourceTs struct {
	Key       string `json:"key,omitempty"`
	Ts 		  int64  `json:"ts,omitempty"`
	Version   int64  `json:"version,omitempty"`
}
