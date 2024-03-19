package bandwidthmanager

import (
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
func (bm *BandwidthManager) RegisterDownload(size int, with *BandwidthManager, tunnelWidth int, tunnelLength int, onDone func(id int)) int {
	newTask := NewTask(size, with, tunnelWidth, tunnelLength, onDone)

	go utils.WithLockedNoResult(&bm.downloadTasksLock, func() {
		bm.downloadTasks = append(bm.downloadTasks, newTask)
		bm.scheduler.Schedule(0)
	})

	return int(newTask.id)
}

func (bm *BandwidthManager) DropDownload(id int) {
	go utils.WithLockedNoResult(&bm.downloadTasksLock, func() {
		for i := 0; i < len(bm.downloadTasks); i++ {
			if bm.downloadTasks[i].id == int64(id) {
				bm.downloadTasks[i] = bm.downloadTasks[len(bm.downloadTasks)-1]
				bm.downloadTasks = bm.downloadTasks[:len(bm.downloadTasks)-1]
				return
			}
		}
		bm.scheduler.Schedule(0)
	})
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

				go task.onDone(int(task.id))
				return false
			}

			return true
		})

		// Adjust speeds
		for i, task := range bm.downloadTasks {
			if task.size == 0 {
				bm.downloadTasks[i].SetSpeed(0)
				continue
			}
			if bm.AvailableDownload() == 0 {
				continue
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
