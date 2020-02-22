package resource

import (
	"encoding/binary"
	"encoding/json"
)

func MarshalResource(r *Resource) ([]byte, error) {
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
	rs := make([]*Resource, count)
	for i := 1; i < count; i++ {
		rs[i] = kind.NewEmptyResource()
	}

	err := json.Unmarshal(data[8:], rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

