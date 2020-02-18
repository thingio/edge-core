package sys

import (
	"bytes"
	"github.com/juju/errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

func UnameKernel() (string, error) {
	data, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		return "", err
	}
	content := string(data)
	parts := strings.Split(content, " ")
	if len(parts) >= 3 {
		return parts[2], nil
	} else {
		return "", errors.Errorf("invalid version string: %s", content)
	}
}

func Date() (string, error) {
	return Exec(`date "+%Y-%m-%d %H:%M:%S %z"`)
}

func Exec(command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", errors.Annotate(err, strings.Trim(stderr.String(), "\n"))

	}
	return stdout.String(), nil
}
