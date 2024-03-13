package bandwidthmanager

import (
	"time"
)

const MaxInt = int(^uint(0) >> 1)

type _Task struct {
	// The size of file in absolute units
	size int
	// The size of dowloaded in absolute units
	done int
	// Absolute units / ms
	workingSpeed int
	// Max Absolute units / ms
	tunnelWidth int
	// Time to start communication in ms
	tunnelLength int
	startedAt    *time.Time

	with      *BandwidthManager
	updatedAt time.Time
	onDone    func()
}

func (t _Task) HasReachedTheOtherSide() bool {
	if t.startedAt == nil {
		return false
	}

	return time.Since(*t.startedAt) >= time.Duration(t.tunnelLength)*time.Millisecond
}

func (t _Task) IsDone() bool {
	return t.done >= t.size
}

func (t *_Task) SetSpeed(absPerMs int) {
	if t.startedAt == nil {
		now := time.Now()
		t.startedAt = &now
	}

	t.workingSpeed = absPerMs
}

// Returns not ok, if download speed at zero
func (t _Task) MsUntilDone() (int, bool) {
	if t.workingSpeed == 0 {
		return MaxInt, false
	}

	if t.HasReachedTheOtherSide() {
		return max((t.size-t.done)/t.workingSpeed, 0), true
	} else {
		timeLeftToReach := time.Duration(t.tunnelLength)*time.Millisecond - time.Since(*t.startedAt)
		return int(timeLeftToReach) + t.size/t.workingSpeed, true
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
		t.done += t.workingSpeed * int(timeSpent.Milliseconds())
	}
	t.updatedAt = time.Now()
}
