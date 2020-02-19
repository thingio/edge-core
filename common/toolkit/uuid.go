package toolkit

import "github.com/rs/xid"

func NewUUID() string {
	return xid.New().String()
}
