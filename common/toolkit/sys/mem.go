package sys

import (
	"github.com/cloudfoundry/gosigar"
	"github.com/juju/errors"
)

// Mem Mem
type Mem struct {
	Total string
	Free  string
}

// GetMem gets memory information
func GetMem() (*Mem, error) {
	m := sigar.Mem{}
	err := m.Get()
	if err != nil {
		return nil, errors.Trace(err)
	}
	mem := &Mem{}
	mem.Total = sigar.FormatSize(m.Total)
	mem.Free = sigar.FormatSize(m.Free)
	return mem, nil
}

// Swap Swap
type Swap struct {
	Total string
	Free  string
}

// GetSwap gets swap space information
func GetSwap() (*Swap, error) {
	s := sigar.Swap{}
	err := s.Get()
	if err != nil {
		return nil, errors.Trace(err)
	}
	swap := &Swap{}
	swap.Total = sigar.FormatSize(s.Total)
	swap.Free = sigar.FormatSize(s.Free)
	return swap, nil
}
