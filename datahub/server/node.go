package server

import (
	"fmt"
	"github.com/alecthomas/units"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/resource"
	"github.com/thingio/edge-core/common/toolkit"
	"github.com/thingio/edge-core/common/toolkit/sys"
	"runtime"
	"strings"
	"syscall"
)

func NodeData(nodeId string) *resource.Node {

	node := new(resource.Node)
	node.Id = nodeId

	node.Stats = map[string]string{}

	var mem, memUsage string
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err == nil {
		mem = fmt.Sprintf("%.2fGiB", float64(sysInfo.Totalram)/float64(units.GiB))
		memUsage = fmt.Sprintf("%.2f%%", float64(sysInfo.Totalram-sysInfo.Freeram)*100.0/float64(sysInfo.Totalram))
	} else {
		log.WithError(err).Error("Fail to get mem usage")
	}

	var cpuUsage string
	if cpu, err := sys.GetCpu(); err == nil {
		cpuUsage = fmt.Sprintf("%.2f%%", float64(cpu.Total-cpu.Idle)*100.0/float64(cpu.Total))
	} else {
		log.WithError(err).Error("Fail to get cpu usage")
	}

	var osVersion string
	if os, err := sys.GetOSVersion(); err == nil {
		osVersion = os
	} else {
		log.WithError(err).Error("Fail to get os version")
	}

	kernelVersion, _ := sys.UnameKernel()
	gpu, _ := sys.GetGpu()
	localtime, _ := sys.Date()

	node.Os = osVersion
	node.Kernel = kernelVersion
	node.Arch = runtime.GOARCH

	node.LocalTime = strings.TrimSpace(localtime)
	node.BootTime = toolkit.NowMS() - sysInfo.Uptime*1e3

	node.Stats["mem"] = mem
	node.Stats["mem_usage"] = memUsage
	node.Stats["cpu"] = fmt.Sprint(runtime.NumCPU())
	node.Stats["cpu_usage"] = cpuUsage
	node.Stats["gpu"] = fmt.Sprint(gpu)

	return node
}
