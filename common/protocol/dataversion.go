package protocol

import (
	"sync/atomic"
	"time"
)

//DataVersion DataVersion
type DataVersion struct {
	timestamp int64
	counter   int64
}

//NextVersion atomic set next version
func (dv *DataVersion) NextVersion() {
	dv.timestamp = time.Now().Unix()
	atomic.AddInt64(&dv.counter, 1)
}

//Timestamp timestamp getter
func (dv *DataVersion) Timestamp() int64 {
	return dv.timestamp
}

//Counter counter getter
func (dv *DataVersion) Counter() int64 {
	return dv.counter
}
