package nlogger

import (
	"dirs/simulation/pkg/node"
	"sync"
	"time"
)

func countMapArr(store map[node.INode][]node.INode) int {
	ans := 0

	for _, arr := range store {
		ans += len(arr)
	}

	return ans
}

func countMapArrWithLock(store map[node.INode][]node.INode, lock *sync.Mutex) int {
	lock.Lock()
	defer lock.Unlock()
	return countMapArr(store)
}

func countMapMapArrWithLock(store map[int]map[node.INode][]node.INode, lock *sync.Mutex) int {
	lock.Lock()
	defer lock.Unlock()

	ans := 0
	for _, emap := range store {
		ans += countMapArr(emap)
	}

	return ans
}

func (l *Logger) CountRouteMessageReceives() int {
	return countMapMapArrWithLock(l.routeMessageReceives, &l.rmrLock)
}
func (l *Logger) CountRouteMessageTimeouts() int {
	return countMapMapArrWithLock(l.routeMessageTimeouts, &l.rmtLock)
}
func (l *Logger) CountRouteMessageConfirms() int {
	return countMapMapArrWithLock(l.routeMessageConfirms, &l.rmcLock)
}
func (l *Logger) CountDeclinedRouteMessageReceives() int {
	return countMapMapArrWithLock(l.deniedRouteMessageReceives, &l.rmrdLock)
}
func (l *Logger) CountDeclinedRouteMessageTimeouts() int {
	return countMapMapArrWithLock(l.deniedRouteMessageTimeouts, &l.rmtdLock)
}
func (l *Logger) CountDeclinedRouteMessageConfirms() int {
	return countMapMapArrWithLock(l.deniedRouteMessageConfirms, &l.rmcdLock)
}

func (l *Logger) CountDownloadMessages() int {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	ans := 0
	for _, earr := range l.downloadMessages {
		ans += len(earr)
	}

	return ans
}

func (l *Logger) DurationToArriveLocked(id int) (time.Duration, bool) {
	l.setLock.Lock()
	defer l.setLock.Unlock()

	arr, ok := l.seTimestamps[id]
	if !ok {
		panic("NO TIMESTAMP FOR ID")
	}

	if arr[1].IsZero() {
		return 0, false
	}

	return arr[1].Sub(arr[0]), true
}

func (l *Logger) AverageDurationToArriveLocked() (time.Duration, int) {
	l.setLock.Lock()
	defer l.setLock.Unlock()

	var durSum time.Duration = 0
	didntReach := 0

	for _, arr := range l.seTimestamps {
		if arr[1].IsZero() {
			didntReach += 1
		} else {
			durSum += arr[1].Sub(arr[0])
		}
	}

	return durSum / time.Duration(len(l.seTimestamps)), didntReach
}

func (l *Logger) AverageRouteMessageReceives() float64 {
	return float64(l.CountRouteMessageReceives()) / float64(len(l.routeMessageReceives))
}
func (l *Logger) AverageRouteMessageTimeouts() float64 {
	return float64(l.CountRouteMessageTimeouts()) / float64(len(l.routeMessageTimeouts))
}
func (l *Logger) AverageRouteMessageConfirms() float64 {
	return float64(l.CountRouteMessageConfirms()) / float64(len(l.routeMessageConfirms))
}
func (l *Logger) AverageDeclinedRouteMessageReceives() float64 {
	return float64(l.CountDeclinedRouteMessageReceives()) / float64(len(l.deniedRouteMessageReceives))
}
func (l *Logger) AverageDeclinedRouteMessageTimeouts() float64 {
	return float64(l.CountDeclinedRouteMessageTimeouts()) / float64(len(l.deniedRouteMessageTimeouts))
}
func (l *Logger) AverageDeclinedRouteMessageConfirms() float64 {
	return float64(l.CountDeclinedRouteMessageConfirms()) / float64(len(l.deniedRouteMessageConfirms))
}

func (l *Logger) AverageDownloadMessages() float64 {
	return float64(l.CountDownloadMessages()) / float64(len(l.downloadMessages))
}
