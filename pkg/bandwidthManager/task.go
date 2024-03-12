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

	with      *BandwidthManager
	updatedAt time.Time
	onDone    func()
}

func (t _Task) IsDone() bool {
	return t.done >= t.size
}

// Returns not ok, if download speed at zero
func (t _Task) MsUntilDone() (int, bool) {
	if t.workingSpeed == 0 {
		return MaxInt, false
	}
	return max((t.size-t.done)/t.workingSpeed, 0), true
}

func (t *_Task) UpdateProgress() {
	timeSpent := time.Since(t.updatedAt)
	t.done += t.workingSpeed * int(timeSpent.Milliseconds())
	t.updatedAt = time.Now()
}
