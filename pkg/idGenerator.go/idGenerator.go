package idgenerator

import "sync/atomic"

var id int32 = 0

func GetId() int {
	defer func() {
		atomic.StoreInt32(&id, atomic.LoadInt32(&id)+1)
	}()

	return int(atomic.LoadInt32(&id))
}
