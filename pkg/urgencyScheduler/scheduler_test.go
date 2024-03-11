package urgencyscheduler

import (
	"testing"
	"time"
)

func TestUrgencyScheduler(t *testing.T) {
	t.Run("Simple test", func(t *testing.T) {
		startTime := time.Now()
		expectToRun := startTime.Add(100 * time.Millisecond)

		s := NewUrgencyScheduler(func() {
			if (time.Since(expectToRun)).Abs() > time.Millisecond*2 {
				t.Errorf("\nRun after %v\nExpected to run after %v", time.Since(startTime), expectToRun.Sub(startTime))
			}
		})

		s.Schedule(100 * time.Millisecond)

		time.Sleep(300 * time.Millisecond)
	})

	t.Run("Reschedule test", func(t *testing.T) {
		startTime := time.Now()
		expectToRun := startTime.Add(200 * time.Millisecond)

		s := NewUrgencyScheduler(func() {
			if (time.Since(expectToRun)).Abs() > time.Millisecond*2 {
				t.Errorf("\nRun after %v\nExpected to run after %v", time.Since(startTime), expectToRun.Sub(startTime))
			}
		})

		s.Schedule(300 * time.Millisecond)

		time.Sleep(10 * time.Millisecond)
		s.Schedule(190 * time.Millisecond)

		time.Sleep(10 * time.Millisecond)
		s.Schedule(300 * time.Millisecond)

		time.Sleep(400 * time.Millisecond)
	})
}
