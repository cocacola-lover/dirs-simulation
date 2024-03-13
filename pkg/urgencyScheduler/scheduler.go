package urgencyscheduler

import (
	"sync/atomic"
	"time"
)

// UrgencyScheduler does not work with the last time given
// but rather with the soonest time given
type UrgencyScheduler struct {
	timedFu        func()
	activationChan chan int32
	innerTimer     int32
}

func (s *UrgencyScheduler) InnerTimer() int32 {
	return atomic.LoadInt32(&s.innerTimer)
}
func (s *UrgencyScheduler) AddInnerTime() {
	atomic.StoreInt32(&s.innerTimer, s.InnerTimer()+1)
}

func (s *UrgencyScheduler) _Watch() {
	for {
		timer, ok := <-s.activationChan
		if !ok {
			return
		}

		if s.InnerTimer() == timer {
			s.AddInnerTime()
			go s.timedFu()
		}
	}
}

func (s *UrgencyScheduler) Schedule(duration time.Duration) {
	go (func() {
		innerTimer := s.InnerTimer()
		time.Sleep(duration)
		s.activationChan <- innerTimer
	})()
}

func (s *UrgencyScheduler) Close() {
	close(s.activationChan)
}

func NewUrgencyScheduler(timedFu func()) *UrgencyScheduler {
	ans := UrgencyScheduler{timedFu: timedFu, activationChan: make(chan int32)}

	go ans._Watch()

	return &ans
}
