package stderr

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type StdError struct {
	Code     int    `json:"code,omitempty"`
	CodeDesc string `json:"code_desc,omitempty"`
	Message  string `json:"message,omitempty"`
	Trace    string `json:"trace,omitempty"`
}

type errCode struct {
	Code       int
	CodeDesc   string
	CodeDescCN string
}

var (
	ErrCodes = make(map[int]*errCode)

	/* General Errors */
	OK             = newErrCode(200, "OK", "成功")
	NotSupported   = newErrCode(300, "Not Supported", "系统不支持")
	NotFound       = newErrCode(400, "Not Found", "数据不存在")
	NotImplemented = newErrCode(600, "Not Implemented", "逻辑未实现")
	Unexpected     = newErrCode(999, "Unexpected Error", "未知错误")

	/* Resource Errors: 100,000 - 100,999 */
	ResourceNotFound      = newErrCode(100404, "Resource Not Found", "资源不存在")
	ResourceAccessFailure = newErrCode(100501, "Resource Access Failure", "资源访问失败")
	ResourceUpdateFailure = newErrCode(100502, "Resource Update Failure", "资源修改失败")
	ResourceDeleteFailure = newErrCode(100503, "Resource Delete Failure", "资源删除失败")

	/* Control Errors: 200,000 - 200,999 */
	InvalidStatusTransition = newErrCode(200501, "Invalid Status Transition", "非法的状态切换")
)

func newErrCode(code int, codeDesc string, codeDescCN string) *errCode {
	if _, ok := ErrCodes[code]; ok {
		log.Panicf("Duplicated error code: %d", code)
	}
	result := &errCode{code, codeDesc, codeDescCN}
	ErrCodes[code] = result
	return result
}

// github.com/juju/errors is a good error library but doesn't support user-defined error code
type locationer interface {
	Location() (string, int)
}

func traceOf(err interface{}) string {
	if err, ok := err.(locationer); ok {
		file, line := err.Location()
		return fmt.Sprintf("%s [%d]", trimGoPath(file), line)
	} else {
		_, file, line, _ := runtime.Caller(2)
		return fmt.Sprintf("%s [%d]", trimGoPath(file), line)
	}
}

var goPath = filepath.Join(build.Default.GOPATH, "src") + string(os.PathSeparator)

func trimGoPath(filename string) string {
	return strings.TrimPrefix(filename, goPath)
}

func (e *errCode) Error(err error) StdError {
	return StdError{e.Code, e.CodeDesc, err.Error(), traceOf(err)}
}

func (e *errCode) Errorf(err error, format string, v ...interface{}) StdError {
	msg := fmt.Sprintf(format, v) + "\nCause: " + err.Error()
	return StdError{e.Code, e.CodeDesc, msg, traceOf(err)}
}

func (e *errCode) Of(format string, v ...interface{}) StdError {
	msg := fmt.Sprintf(format, v)
	return StdError{e.Code, e.CodeDesc, msg, traceOf(nil)}
}

func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}