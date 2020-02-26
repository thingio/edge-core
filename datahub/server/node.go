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

	var osVersion string
	if os, err := sys.GetOSVersion(); err == nil {
		osVersion = os
	} else {
		log.WithError(err).Error("Fail to get os version")
	}

	node.Os = osVersion
	node.Kernel, _ = sys.UnameKernel()
	node.Arch = runtime.GOARCH

	return node
}

func NodeState(nodeId string) *resource.State {
	node := resource.State{}
	node["id"] = nodeId

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

	gpu, _ := sys.GetGpu()
	localtime, _ := sys.Date()

	node["local_time"] = strings.TrimSpace(localtime)
	node["boot_time"] = fmt.Sprint(toolkit.NowMS() - sysInfo.Uptime*1e3)

	node["mem"] = mem
	node["mem_usage"] = memUsage
	node["cpu"] = fmt.Sprint(runtime.NumCPU())
	node["cpu_usage"] = cpuUsage
	node["gpu"] = fmt.Sprint(gpu)

	return &node
}
