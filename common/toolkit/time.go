package toolkit

import "time"

func NowMS() int64 {
	return time.Now().UnixNano() / 1e6
}
