package resource

import (
	"encoding/binary"
	"encoding/json"
)

func MarshalResource(r *Resource) ([]byte, error) {
	// json.Marshal(nil) will returns 'null', which is not neat
	if r == nil {
		return make([]byte, 0), nil
	}
	return json.Marshal(r)
}

func UnmarshalResource(kind Kind, data []byte) (*Resource, error) {
	r := kind.NewResource()
	err := json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func MarshalResourceList(rs []*Resource) ([]byte, error) {
	// first 8 bytes represents the length of the resource
	if len(rs) == 0 {
		return make([]byte, 8), nil
	}

	data, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}
	dataLen := len(data)
	count := uint64(len(rs))
	result := make([]byte, 8 + dataLen)
	binary.BigEndian.PutUint64(result, count)
	copy(result[8:], data)
	return result, nil
}

func UnmarshalResourceList(kind Kind, data []byte) ([]*Resource, error) {
	count := int(binary.BigEndian.Uint64(data[:8]))
	if count == 0 {
		return nil, nil
	}

	rs := make([]*Resource, count)
	for i := 0; i < count; i++ {
		rs[i] = kind.NewEmptyResource()
	}

	data = data[8:]
	err := json.Unmarshal(data, &rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

