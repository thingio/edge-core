package resource

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	p := Pipeline{
		Id:      "test-123",
		Type:    "ts",
		Name:    "just a test",
		BodyDef: "just a body",
	}
	data, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
		return
	}

	expect := `{"id":"test-123","type":"ts","name":"just a test","body_def":"just a body"}`
	if expect != string(data) {
		t.Fail()
	}

}

func TestUnmarshal(t *testing.T) {
	data := []byte(`{"type":"ts","name":"just a test","body_def":"just a body"}`)
	p := KindPipeline.NewResource()
	id := p.Value.(*Pipeline).GetId()
	fmt.Printf("%+v\n", p.Value)
	err := json.Unmarshal(data, p.Value)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", p.Value)
	if id != p.Value.(*Pipeline).GetId() {
		t.Fail()
	}
}

func TestUnmarshalList(t *testing.T) {
	ps := []*Resource{KindPipeline.NewResource(), KindPipeline.NewResource()}
	data, err := json.Marshal(ps)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v %+v\n", ps[0].Value, ps[1].Value)

	newPs := []*Resource{}
	size := 5
	for i := 1; i <= size; i++ {
		newPs = append(newPs, KindPipeline.NewResource())
	}

	err = json.Unmarshal(data, &newPs)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v %+v\n", newPs[0].Value, newPs[1].Value)
}

func TestUnmarshalNull(t *testing.T) {
	data, _ := json.Marshal(nil)
	fmt.Println("marshal", nil, string(data))

	var value interface{}
	_ = json.Unmarshal(data, &value)
	fmt.Println("unmarshal", data, value)

	if []byte("") == nil {
		fmt.Println("yes")
	}

	a := make([]byte, 0)
	fmt.Println(a)
}

func TestUnmarshalDevice(t *testing.T) {
	data := []byte(`[{"node_id":"my-node","kind":"device","id":"device-bp970lpsem6vg12t1gug","value":{"id":"string","name":"string","product":"string","props":{"additionalProp1":"string","additionalProp2":"string","additionalProp3":"string"}},"ts":1582461015635431162,"version":1}]`)

	value := make([]*Resource, 1)
	value[0] = KindDevice.NewEmptyResource()
	err := json.Unmarshal(data, &value)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", value)
}

func TestUnmarshalDeviceList(t *testing.T) {
	data := []byte(`[{"node_id":"my-node","kind":"device","id":"device-bp970lpsem6vg12t1gug","value":{"id":"string","name":"string","product":"string","props":{"additionalProp1":"string","additionalProp2":"string","additionalProp3":"string"}},"ts":1582461015635431162,"version":1}]`)
	dataLen := len(data)
	result := make([]byte, 8+dataLen)
	binary.BigEndian.PutUint64(result, 1)
	copy(result[8:], data)

	value, err := UnmarshalResourceList(KindDevice, result)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", value)
}
