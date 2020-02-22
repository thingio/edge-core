package toolkit

import "time"

func Now() int64 {
	return time.Now().UnixNano()
}

func NowMS() int64 {
	return time.Now().UnixNano() / 1e6
}
