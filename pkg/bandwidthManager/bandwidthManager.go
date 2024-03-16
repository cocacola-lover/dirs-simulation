package bandwidthmanager

import (
	netp "dirs/simulation/pkg/network"
	usp "dirs/simulation/pkg/urgencyScheduler"
	"dirs/simulation/pkg/utils"
	"sync"
)

type BandwidthManager struct {
	maxUpload    int
	maxDownload  int
	usedUpload   int32
	usedDownload int32

	downloadTasks     []_Task
	downloadTasksLock sync.Mutex

	scheduler *usp.UrgencyScheduler
}

// Use with go
func (bm *BandwidthManager) RegisterDownload(size int, with *BandwidthManager, tunnel netp.Tunnel, onDone func()) {
	utils.WithLockedNoResult(&bm.downloadTasksLock, func() {
		bm.downloadTasks = append(bm.downloadTasks, NewTask(size, with, tunnel, onDone))
	})

	bm.scheduler.Schedule(0)
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

			if canMake := min(
				task.workingSpeed+bm.AvailableDownload(),
				task.with.AvailableUpload(),
				task.tunnelWidth); canMake > task.workingSpeed {

				bm.downloadTasks[i].SetSpeed(canMake)

				bm._AddDownload(canMake - task.workingSpeed)
				task.with._AddUpload(canMake - task.workingSpeed)
			}
		}

		// Schedule another reevaluation
		if len(bm.downloadTasks) == 0 {
			return
		}

		soonestToEnd, _ := bm.downloadTasks[0].UntilDone()
		for i := 1; i < len(bm.downloadTasks); i++ {
			newValue, _ := bm.downloadTasks[i].UntilDone()
			soonestToEnd = min(newValue, soonestToEnd)
		}

		bm._ScheduleReevaluation(soonestToEnd)
	})
}

func NewBandwidthManager(maxDownload, maxUpload int) *BandwidthManager {
	bm := &BandwidthManager{maxUpload: maxUpload, maxDownload: maxDownload}

	bm.scheduler = usp.NewUrgencyScheduler(bm._Reevaluate)

	return bm
}