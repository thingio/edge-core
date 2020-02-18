package toolkit

import (
	"fmt"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto"
	"os"
	"runtime"
	"strings"
	"syscall"
	"github.com/alecthomas/units"
)

func GetEdgeNodeInfo(edgeId, edgeHost string) *proto.Node {
	edge := new(proto.Node)
	edge.Id = edgeId
	edge.Name = edgeId
	if edgeHost != "" {
		edge.Host = edgeHost
	} else {
		edge.Host, _ = os.Hostname()
	}

	edge.Infos = map[string]string{}

	var mem, memUsage string
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err == nil {
		mem = fmt.Sprintf("%.2fGiB", float64(sysInfo.Totalram)/float64(units.GiB))
		memUsage = fmt.Sprintf("%.2f%%", float64(sysInfo.Totalram-sysInfo.Freeram)*100.0/float64(sysInfo.Totalram))
	} else {
		log.WithError(err).Error("STDHUB: fail to get mem usage")
	}

	var cpuUsage string
	if cpu, err := sys.GetCpu(); err == nil {
		cpuUsage = fmt.Sprintf("%.2f%%", float64(cpu.Total-cpu.Idle)*100.0/float64(cpu.Total))
	} else {
		log.WithError(err).Error("STDHUB: fail to get cpu usage")
	}

	var osVersion string
	if os, err := sys.GetOSVersion(); err == nil {
		osVersion = os
	} else {
		log.WithError(err).Error("STDHUB: fail to get os version")
	}

	kernelVersion, _ := sys.UnameKernel()
	gpu, _ := sys.GetGpu()
	localtime, _ := sys.Date()

	edge.Infos["localtime"] = strings.TrimSpace(localtime)
	edge.Infos["mem"] = mem
	edge.Infos["mem_usage"] = memUsage
	edge.Infos["cpu"] = fmt.Sprint(runtime.NumCPU())
	edge.Infos["cpu_usage"] = cpuUsage
	edge.Infos["gpu"] = fmt.Sprint(gpu)

	edge.Infos["os"] = osVersion
	edge.Infos["kernel"] = kernelVersion
	edge.Infos["arch"] = runtime.GOARCH
	edge.CreateTime = toolkit.NowMS() - sysInfo.Uptime*1e3
	edge.HeartbeatTime = toolkit.NowMS()
	return edge
}

