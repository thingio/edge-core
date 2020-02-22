package resource

import (
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
	fmt.Println(string(data))

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
