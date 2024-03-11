package bandwidthmanager

import (
	"time"
)

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

func (t _Task) MsUntilDone() int {
	return max((t.size-t.done)/t.workingSpeed, 0)
}

func (t *_Task) UpdateProgress() {
	timeSpent := time.Since(t.updatedAt)
	t.done += t.workingSpeed * int(timeSpent.Milliseconds())
	t.updatedAt = time.Now()
}
