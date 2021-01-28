package protocol

import (
	"sync/atomic"
	"time"
)

type DataVersion struct {
	timestamp int64
	counter   int64
}

func (dv *DataVersion) NextVersion() {
	dv.timestamp = time.Now().Unix()
	atomic.AddInt64(&dv.counter, 1)
}

func (dv *DataVersion) Timestamp() int64 {
	return dv.timestamp
}

func (dv *DataVersion) Counter() int64 {
	return dv.counter
}
