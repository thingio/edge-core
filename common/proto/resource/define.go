package resource

import (
	"fmt"
)

const (
	GenusTS = "ts" // 时序分类
	GenusMM = "mm" // 多媒体分类
)

type IdObject interface {
	GetId() string
	SetId(id string)
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

const StateKeyStatus = "status"

const (
	StatusInit = "INIT"
	StatusUp   = "UP"
	StatusKill = "KILL"
	StatusDown = "DOWN"
	StatusFail = "FAIL"
)
