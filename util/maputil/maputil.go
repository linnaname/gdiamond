package maputil

import (
	"sync"
	"sync/atomic"
)

//LengthOfSyncMap get length of sync.Map,I can't figure out a better way,not thread-safe
func LengthOfSyncMap(sm sync.Map) int64 {
	l := int64(0)
	sm.Range(func(k, v interface{}) bool {
		atomic.AddInt64(&l, 1)
		return true
	})
	return l
}

//ClearSyncMap  clear all k,v of sync.Map,not thread-safe
func ClearSyncMap(sm sync.Map) {
	sm.Range(func(k, v interface{}) bool {
		sm.Delete(k)
		return true
	})
}
