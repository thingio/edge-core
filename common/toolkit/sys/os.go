package sys

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/host"
)

func GetOSVersion() (string, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return "", errors.Errorf("fail to fetch os version: %s", err)
	}
	return fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion), nil
}
