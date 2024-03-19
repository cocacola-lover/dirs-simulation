package bandwidthmanager

import (
	"testing"
	"time"
)

func TestBandwidthManager(t *testing.T) {
	t.Run("Simple test", func(t *testing.T) {
		startTime := time.Now()
		expectToRun := startTime.Add(30 * time.Millisecond)

		bm1 := NewBandwidthManager(10, 10)
		bm2 := NewBandwidthManager(10, 10)

		bm1.RegisterDownload(100, bm2, 5, 10, func(_ int) {
			if time.Since(expectToRun).Abs() > 2*time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime), expectToRun.Sub(startTime))
			}
			if bm1.scheduler.InnerTimer() != 2 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 2)
			}
		})

		time.Sleep(30 * time.Millisecond)
	})

	t.Run("Simple test, 1 speed", func(t *testing.T) {
		startTime := time.Now()
		expectToRun := startTime.Add(101 * time.Millisecond)

		bm1 := NewBandwidthManager(1, 1)
		bm2 := NewBandwidthManager(1, 1)

		bm1.RegisterDownload(100, bm2, 1, 1, func(_ int) {
			if time.Since(expectToRun).Abs() > 2*time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime), expectToRun.Sub(startTime))
			}
			if bm1.scheduler.InnerTimer() != 2 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 2)
			}
		})

		time.Sleep(110 * time.Millisecond)
	})

	t.Run("Simple test, overflow tasks", func(t *testing.T) {
		startTime := time.Now()
		expectToRun1 := startTime.Add(10 * time.Millisecond)
		expectToRun2 := startTime.Add(15 * time.Millisecond)

		bm1 := NewBandwidthManager(1, 1)
		bm2 := NewBandwidthManager(1, 1)

		bm1.RegisterDownload(10, bm2, 1, 0, func(_ int) {
			if time.Since(expectToRun1).Abs() > 2*time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime), expectToRun1.Sub(startTime))
			}
			if bm1.scheduler.InnerTimer() > 3 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 3)
			}
		})

		time.Sleep(time.Millisecond)

		bm1.RegisterDownload(5, bm2, 1, 0, func(_ int) {
			if time.Since(expectToRun2).Abs() > 2*time.Millisecond {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime), expectToRun2.Sub(startTime))
			}
			if bm1.scheduler.InnerTimer() > 4 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 4)
			}
		})

		time.Sleep(110 * time.Millisecond)
	})

	t.Run("Dual test", func(t *testing.T) {
		startTime1 := time.Now()
		expectToRun1 := startTime1.Add(10 * time.Millisecond)

		bm1 := NewBandwidthManager(10, 10)
		bm2 := NewBandwidthManager(10, 10)
		bm3 := NewBandwidthManager(10, 10)

		bm1.RegisterDownload(70, bm2, 7, 0, func(_ int) {
			if time.Since(expectToRun1).Abs() > time.Millisecond*2 {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime1), expectToRun1.Sub(startTime1))
			}
			if bm1.scheduler.InnerTimer() != 3 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 3)
			}
		})

		time.Sleep(5 * time.Millisecond)

		startTime2 := time.Now()
		expectToRun2 := startTime2.Add(11 * time.Millisecond)

		bm1.RegisterDownload(45, bm3, 5, 0, func(_ int) {
			if time.Since(expectToRun2).Abs() > time.Millisecond*2 {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime2), expectToRun2.Sub(startTime2))
			}
			if bm1.scheduler.InnerTimer() != 4 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 4)
			}
		})

		time.Sleep(30 * time.Millisecond)
	})

	t.Run("Cross download test", func(t *testing.T) {
		startTime1 := time.Now()
		expectToRun1 := startTime1.Add(10 * time.Millisecond)

		bm1 := NewBandwidthManager(10, 10)
		bm2 := NewBandwidthManager(10, 10)

		bm1.RegisterDownload(90, bm2, 12, 1, func(_ int) {
			if time.Since(expectToRun1).Abs() > time.Millisecond*5 {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime1), expectToRun1.Sub(startTime1))
			}
			if bm1.scheduler.InnerTimer() > 2 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 2)
			}
		})

		time.Sleep(5 * time.Millisecond)

		startTime2 := time.Now()
		expectToRun2 := startTime2.Add(9 * time.Millisecond)

		bm2.RegisterDownload(45, bm1, 5, 0, func(_ int) {
			if time.Since(expectToRun2).Abs() > time.Millisecond*5 {
				t.Errorf("Upload took %v, but was expected to take %v", time.Since(startTime2), expectToRun2.Sub(startTime2))
			}
			if bm1.scheduler.InnerTimer() > 2 {
				t.Errorf("Upload took %v reevaluations, but was expected to take %v", bm1.scheduler.InnerTimer(), 2)
			}
		})

		time.Sleep(30 * time.Millisecond)
	})

	t.Run("Test DropDownload", func(t *testing.T) {
		bm1 := NewBandwidthManager(10, 10)
		bm2 := NewBandwidthManager(10, 10)

		downloadId := bm1.RegisterDownload(10, bm2, 5, 10, func(_ int) {
			t.Fatal("Download did not drop")
		})

		time.Sleep(5 * time.Millisecond)

		bm1.DropDownload(downloadId)

		time.Sleep(30 * time.Millisecond)
	})
}
