package bandwidthmanager

import (
	"sync/atomic"
	"time"
)

func (bm *BandwidthManager) AvailableDownload() int {
	return bm.maxDownload - int(atomic.LoadInt32(&bm.usedDownload))
}

func (bm *BandwidthManager) AvailableUpload() int {
	return bm.maxUpload - int(atomic.LoadInt32(&bm.usedUpload))
}

func (bm *BandwidthManager) _AddDownload(val int) {
	usedDownload := atomic.LoadInt32(&bm.usedDownload)
	atomic.StoreInt32(&bm.usedDownload, usedDownload+int32(val))
}

func (bm *BandwidthManager) _AddUpload(val int) {
	usedUpload := atomic.LoadInt32(&bm.usedUpload)
	atomic.StoreInt32(&bm.usedUpload, usedUpload+int32(val))
}

func (bm *BandwidthManager) _ScheduleReevaluation(dur time.Duration) {
	bm.scheduler.Schedule(dur)
}
