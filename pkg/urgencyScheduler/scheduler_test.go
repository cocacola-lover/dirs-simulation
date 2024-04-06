package urgencyscheduler

import (
	"runtime"
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

func TestUrgencyScheduler_Close(t *testing.T) {
	norm := runtime.NumGoroutine()

	s := NewUrgencyScheduler(func() {})

	if runtime.NumGoroutine() != norm+1 {
		t.Errorf("Expected to run %d goroutines but encountered %d\n", norm+1, runtime.NumGoroutine())
	}

	s.Schedule(time.Millisecond * 10)

	if runtime.NumGoroutine() != norm+2 {
		t.Errorf("Expected to run %d goroutines but encountered %d\n", norm+2, runtime.NumGoroutine())
	}

	time.Sleep(time.Millisecond * 10)

	if runtime.NumGoroutine() != norm+1 {
		t.Errorf("Expected to run %d goroutines but encountered %d\n", norm+1, runtime.NumGoroutine())
	}

	s.Schedule(time.Millisecond * 10)
	s.Schedule(time.Millisecond * 5)

	if runtime.NumGoroutine() != norm+3 {
		t.Errorf("Expected to run %d goroutines but encountered %d\n", norm+3, runtime.NumGoroutine())
	}

	s.Close()

	time.Sleep(time.Millisecond * 10)

	if runtime.NumGoroutine() != norm {
		t.Errorf("Expected to run %d goroutines but encountered %d\n", norm, runtime.NumGoroutine())
	}
}
