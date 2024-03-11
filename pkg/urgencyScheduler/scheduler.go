package urgencyscheduler

import (
	"dirs/simulation/pkg/utils"
	"sync"
	"time"
)

// UrgencyScheduler does not work with the last time given
// but rather with the soonest time given
type UrgencyScheduler struct {
	fu             func()
	activationChan chan int

	innerTimer     int
	innerTimerLock sync.Mutex
}

func (s *UrgencyScheduler) _Watch() {
	for {
		timer, ok := <-s.activationChan
		if !ok {
			return
		}

		utils.WithLockedNoResult(&s.innerTimerLock, func() {
			if s.innerTimer == timer {
				s.innerTimer++
				go s.fu()
			}
		})
	}
}

func (s *UrgencyScheduler) Schedule(duration time.Duration) {
	go (func() {
		innerTimer := utils.WithLocked(&s.innerTimerLock, func() int {
			return s.innerTimer
		})

		time.Sleep(duration)
		s.activationChan <- innerTimer
	})()
}

func (s *UrgencyScheduler) Close() {
	close(s.activationChan)
}

func NewUrgencyScheduler(fu func()) *UrgencyScheduler {
	ans := UrgencyScheduler{fu: fu, activationChan: make(chan int)}

	go ans._Watch()

	return &ans
}
