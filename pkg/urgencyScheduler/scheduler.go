package urgencyscheduler

import (
	"sync/atomic"
	"time"
)

// UrgencyScheduler does not work with the last time given
// but rather with the soonest time given
type UrgencyScheduler struct {
	timedFu    func()
	innerTimer int32

	activationChan  chan int32
	timeToCloseChan chan bool
}

func (s *UrgencyScheduler) InnerTimer() int32 {
	return atomic.LoadInt32(&s.innerTimer)
}
func (s *UrgencyScheduler) AddInnerTime() {
	atomic.StoreInt32(&s.innerTimer, s.InnerTimer()+1)
}

func (s *UrgencyScheduler) _Watch() {
	for {
		select {
		case timer := <-s.activationChan:
			if s.InnerTimer() == timer {
				s.AddInnerTime()
				go s.timedFu()
			}
		case <-s.timeToCloseChan:
			return
		}
	}
}

func (s *UrgencyScheduler) Schedule(duration time.Duration) {
	go (func() {
		innerTimer := s.InnerTimer()
		time.Sleep(duration)

		if !s.IsClosed() {
			s.activationChan <- innerTimer
		}
	})()
}

func (s *UrgencyScheduler) IsClosed() bool {
	select {
	case <-s.timeToCloseChan:
		return true
	default:
		return false
	}
}

func (s *UrgencyScheduler) Close() {
	close(s.timeToCloseChan)
}

func NewUrgencyScheduler(timedFu func()) *UrgencyScheduler {
	ans := UrgencyScheduler{
		timedFu:         timedFu,
		activationChan:  make(chan int32),
		timeToCloseChan: make(chan bool),
	}

	go ans._Watch()

	return &ans
}
