package bandwidthmanager

import (
	"sync/atomic"
	"time"
)

var id int64 = 0

type _Task struct {
	id int64
	// The size of file in absolute units
	size int
	// The size of dowloaded in absolute units
	done float64
	// Absolute units / ms
	workingSpeed int
	// Max Absolute units / ms
	tunnelWidth int
	// Time to start communication in ms
	tunnelLength int
	startedAt    *time.Time

	with      *BandwidthManager
	updatedAt time.Time
	onDone    func(id int)
}

func (t _Task) HasReachedTheOtherSide() bool {
	if t.startedAt == nil {
		return false
	}

	return time.Since(*t.startedAt) >= time.Duration(t.tunnelLength)*time.Millisecond
}

func (t _Task) IsDone() bool {
	if t.startedAt == nil {
		return false
	}
	if time.Since(*t.startedAt).Milliseconds() < int64(t.tunnelLength) {
		return false
	}
	return t.done >= float64(t.size)
}

func (t *_Task) SetSpeed(absPerMs int) {
	if t.startedAt == nil {
		now := time.Now()
		t.startedAt = &now
	}

	t.workingSpeed = absPerMs
}

// Returns not ok, if download speed at zero
func (t _Task) UntilDone() (time.Duration, bool) {
	if t.workingSpeed == 0 {
		if t.size != 0 {
			return time.Millisecond * 10, false
		}
		return time.Duration(t.tunnelLength)*time.Millisecond - time.Since(*t.startedAt), true
	}
	if t.HasReachedTheOtherSide() {
		timeLeftToInstall := time.Duration((float64(t.size) - t.done) * float64(time.Millisecond) / float64(t.workingSpeed))
		return max(timeLeftToInstall, 0), true
	} else {
		timeLeftToReach := time.Duration(t.tunnelLength)*time.Millisecond - time.Since(*t.startedAt)
		timeToInstall := time.Duration(float32(t.size*int(time.Millisecond)) / float32(t.workingSpeed))

		return timeLeftToReach + timeToInstall, true
	}

}

func (t *_Task) UpdateProgress() {
	if t.HasReachedTheOtherSide() {
		var countingFrom time.Time

		if t.updatedAt.After(*t.startedAt) {
			countingFrom = t.updatedAt
		} else {
			countingFrom = *t.startedAt
		}

		timeSpent := time.Since(countingFrom)
		t.done += float64(t.workingSpeed*int(timeSpent)) / float64(time.Millisecond)
	}
	t.updatedAt = time.Now()
}

func NewTask(size int, with *BandwidthManager, tunnelWidth int, tunnelLength int, onDone func(id int)) _Task {
	lid := atomic.LoadInt64(&id)
	defer atomic.StoreInt64(&id, lid+1)

	return _Task{id: lid, size: size, with: with, tunnelWidth: tunnelWidth, tunnelLength: tunnelLength, onDone: onDone, updatedAt: time.Now()}
}
