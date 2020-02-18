package sys

import (
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/pkg/errors"
	"runtime"
	"strconv"
	"strings"
)

// Gpu Gpu
type Gpu struct {
	ID     string
	Model  string
	Memory Mem
}

// GetGpu gets gpu information
func GetGpu() ([]Gpu, error) {
	var gpus []Gpu
	if strings.HasPrefix(runtime.GOARCH, "amd") {
		out, err := Exec(`nvidia-smi --query-gpu=index,name,memory.total,memory.free --format=csv,noheader,nounits`)
		if err != nil {
			return gpus, errors.Errorf("can't get cpu info: %v", err)
		}

		for _, raw := range strings.Split(out, "\n") {
			if strings.TrimSpace(raw) == "" {
				continue
			}
			var g Gpu
			t := strings.Split(raw, ",")
			total, err := strconv.Atoi(strings.TrimSpace(t[2]))
			if err != nil {
				return gpus, err
			}
			free, err := strconv.Atoi(strings.TrimSpace(t[3]))
			if err != nil {
				return gpus, err
			}
			g.ID = strings.TrimSpace(t[0])
			g.Model = strings.TrimSpace(t[1])
			g.Memory = Mem{
				sigar.FormatSize(uint64(total * 1024 * 1024)),
				sigar.FormatSize(uint64(free * 1024 * 1024)),
			}
			gpus = append(gpus, g)
		}
		return gpus, nil
	}
	return nil, errors.Errorf("can't support getting gpu info")
}

func GetGpuType() (string, error) {
	gpus, err := GetGpu()
	if err != nil || len(gpus) == 0 {
		return "", err
	}
	gpuModels := make([]string, 0)
	for _, gpu := range gpus {
		gpuModels = append(gpuModels, gpu.Model)
	}
	return strings.Join(gpuModels, ";"), nil
}
