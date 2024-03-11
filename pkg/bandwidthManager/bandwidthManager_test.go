package bandwidthmanager

import (
	"testing"
	"time"
)

func TestBandwidthManager(t *testing.T) {
	t.Run("Simple test", func(t *testing.T) {
		startTime := time.Now()
		expectToRun := startTime.Add(20 * time.Millisecond)

		bm1 := NewBandwidthManager(10, 10)
		bm2 := NewBandwidthManager(10, 10)

		bm1.RegisterDownload(100, bm2, 5, func() {
			t.Log("Download done")
			if time.Since(expectToRun).Abs() > time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime), expectToRun.Sub(startTime))
			}
			if bm1.scheduler.InnerTimer() != 2 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 2)
			}
		})

		time.Sleep(30 * time.Millisecond)
	})

	t.Run("Dual test", func(t *testing.T) {
		startTime1 := time.Now()
		expectToRun1 := startTime1.Add(10 * time.Millisecond)

		bm1 := NewBandwidthManager(10, 10)
		bm2 := NewBandwidthManager(10, 10)
		bm3 := NewBandwidthManager(10, 10)

		bm1.RegisterDownload(70, bm2, 7, func() {
			t.Log("Download done")
			if time.Since(expectToRun1).Abs() > time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime1), expectToRun1.Sub(startTime1))
			}
			if bm1.scheduler.InnerTimer() != 3 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 3)
			}
		})

		time.Sleep(5 * time.Millisecond)

		startTime2 := time.Now()
		expectToRun2 := startTime2.Add(11 * time.Millisecond)

		bm1.RegisterDownload(45, bm3, 5, func() {
			t.Log("Download done")
			if time.Since(expectToRun2).Abs() > time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime2), expectToRun2.Sub(startTime2))
			}
			if bm1.scheduler.InnerTimer() != 4 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 4)
			}
		})

		time.Sleep(30 * time.Millisecond)
	})
}
