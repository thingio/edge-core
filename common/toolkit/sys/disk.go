package sys

import (
	"github.com/cloudfoundry/gosigar"
	"github.com/juju/errors"
)

// Disk Disk
type Disk struct {
	Total string
	Free  string
}

// GetDisk gets disk information
func GetDisk() (*Disk, error) {
	t := uint64(0)
	f := uint64(0)

	fslist := sigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, fs := range fslist.List {
		usage := sigar.FileSystemUsage{}
		err := usage.Get(fs.DirName)
		if err != nil {
			return nil, errors.Trace(err)
		}
		t = t + usage.Total
		f = f + usage.Avail
	}
	disk := &Disk{}
	disk.Total = sigar.FormatSize(t * 1024)
	disk.Free = sigar.FormatSize(f * 1024)
	return disk, nil
}
