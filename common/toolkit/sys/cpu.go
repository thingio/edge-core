package sys

import (
	"fmt"
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/juju/errors"
	"github.com/shirou/gopsutil/cpu"
)

type Cpu struct {
	User  uint64
	Idle  uint64
	Wait  uint64
	Total uint64
}

var prev *sigar.Cpu
// GetMem gets memory information
func GetCpu() (*Cpu, error) {
	c := new(sigar.Cpu)
	err := c.Get()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if prev == nil {
		prev = c
		return &Cpu{
			c.User,
			c.Idle,
			c.Wait,
			c.Total(),
		}, nil
	} else {
		dc := c.Delta(*prev)
		prev = &dc
		return &Cpu{
			dc.User,
			dc.Idle,
			dc.Wait,
			dc.Total(),
		}, nil
	}
}

func GetCpuType() (string, error) {
	cpuInfos, err := cpu.Info()
	if err != nil || len(cpuInfos) == 0 {
		return "", errors.Errorf("fail to fetch cpu type: %s", err)
	}
	return fmt.Sprintf("%s: %d cores", cpuInfos[0].ModelName, len(cpuInfos)), nil
}
