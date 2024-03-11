package bandwidthmanager

import (
	usp "dirs/simulation/pkg/urgencyScheduler"
	"dirs/simulation/pkg/utils"
	"sync"
	"time"
)

type BandwidthManager struct {
	maxUpload   int
	maxDownload int

	usedUpload     int
	usedUploadLock sync.Mutex

	usedDownload     int
	usedDownloadLock sync.Mutex

	downloadTasks     []_Task
	downloadTasksLock sync.Mutex

	scheduler *usp.UrgencyScheduler
}

func (bm *BandwidthManager) AvailableDownload() int {
	return utils.WithLocked(&bm.usedDownloadLock, func() int {
		return bm.maxDownload - bm.usedDownload
	})
}

func (bm *BandwidthManager) AvailableUpload() int {
	return utils.WithLocked(&bm.usedUploadLock, func() int {
		return bm.maxUpload - bm.usedUpload
	})
}

func (bm *BandwidthManager) _AddDownload(val int) {
	utils.WithLockedNoResult(&bm.usedDownloadLock, func() {
		bm.usedDownload += val
	})
}

func (bm *BandwidthManager) _AddUpload(val int) {
	utils.WithLockedNoResult(&bm.usedUploadLock, func() {
		bm.usedUpload += val
	})
}

// Use with go
func (bm *BandwidthManager) RegisterDownload(size int, with *BandwidthManager, tunnelWidth int, onDone func()) {
	utils.WithLockedNoResult(&bm.downloadTasksLock, func() {
		bm.downloadTasks = append(bm.downloadTasks, _Task{size: size, with: with, tunnelWidth: tunnelWidth, onDone: onDone, updatedAt: time.Now()})
	})

	bm.scheduler.Schedule(0)
}

func (bm *BandwidthManager) _ScheduleReevaluation(inMs int) {
	bm.scheduler.Schedule(time.Millisecond * time.Duration(inMs))
}

func (bm *BandwidthManager) _Reevaluate() {

	utils.WithLockedNoResult(&bm.downloadTasksLock, func() {
		for i := range bm.downloadTasks {
			bm.downloadTasks[i].UpdateProgress()
		}

		// Check if anything has freed up
		bm.downloadTasks = utils.Filter(bm.downloadTasks, func(task _Task, ind int) bool {

			if task.IsDone() {
				bm._AddDownload(-task.workingSpeed)
				task.with._AddUpload(-task.workingSpeed)

				go task.onDone()
				return false
			}

			return true
		})

		// Adjust speeds
		for i, task := range bm.downloadTasks {
			if bm.AvailableDownload() == 0 {
				break
			}

			if canMake := min(task.workingSpeed+bm.AvailableDownload(), task.with.AvailableUpload(), task.tunnelWidth); canMake > task.workingSpeed {
				bm.downloadTasks[i].workingSpeed = canMake

				bm._AddDownload(canMake - task.workingSpeed)
				task.with._AddUpload(canMake - task.workingSpeed)
			}
		}

		// Schedule another reevaluation
		if len(bm.downloadTasks) == 0 {
			return
		}

		soonestToEnd := bm.downloadTasks[0].MsUntilDone()
		for i := 1; i < len(bm.downloadTasks); i++ {
			soonestToEnd = min(bm.downloadTasks[i].MsUntilDone(), soonestToEnd)
		}

		bm._ScheduleReevaluation(soonestToEnd)
	})
}

func NewBandwidthManager(maxDownload, maxUpload int) *BandwidthManager {
	bm := &BandwidthManager{maxUpload: maxUpload, maxDownload: maxDownload}

	bm.scheduler = usp.NewUrgencyScheduler(bm._Reevaluate)

	return bm
}
