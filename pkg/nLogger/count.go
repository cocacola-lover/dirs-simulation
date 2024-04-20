package nlogger

import (
	"dirs/simulation/pkg/node"
	"sync"
	"sync/atomic"
	"time"
)

func countMapArr(store map[node.INode]int) int {
	ans := 0

	for _, v := range store {
		ans += v
	}

	return ans
}

func countMapArrWithLock(store map[node.INode]int, lock *sync.Mutex) int {
	lock.Lock()
	defer lock.Unlock()
	return countMapArr(store)
}

func countMapMapArrWithLock(store map[int]int, lock *sync.Mutex) int {
	lock.Lock()
	defer lock.Unlock()

	ans := 0
	for _, v := range store {
		ans += v
	}

	return ans
}

func (l *Logger) CountRouteMessageReceives() int {
	return countMapMapArrWithLock(l.routeMessageReceives, &l.rmrLock)
}
func (l *Logger) CountFaultMessageReceives() int {
	return int(atomic.LoadInt32(&l.faultMessageReceives))
}
func (l *Logger) CountRouteMessageConfirms() int {
	return countMapMapArrWithLock(l.routeMessageConfirms, &l.rmcLock)
}
func (l *Logger) CountDeclinedRouteMessageReceives() int {
	return countMapMapArrWithLock(l.deniedRouteMessageReceives, &l.rmrdLock)
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

func (l *Logger) countDownloadPathWithoutLock(id int) int {
	lead, ok1 := l.startedSearches[id]
	dict, ok2 := l.downloadMessages[id]
	if !ok1 || !ok2 {
		panic("Counting download path of nonexistent download")
	}

	ans := []node.INode{lead}

	for i := 0; true; i++ {
		next, ok := dict[ans[i]]
		if !ok {
			break
		}

		ans = append(ans, next)
	}

	return len(ans) - 1
}

func (l *Logger) CountDownloadPath(id int) int {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	return l.countDownloadPathWithoutLock(id)
}

func (l *Logger) CountFailedNodes() int {
	return int(atomic.LoadInt32(&l.failedNode))
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
func (l *Logger) AverageFaultMessageReceives() float64 {
	return float64(l.CountFaultMessageReceives()) / float64(len(l.routeMessageReceives))
}
func (l *Logger) AverageRouteMessageConfirms() float64 {
	return float64(l.CountRouteMessageConfirms()) / float64(len(l.routeMessageConfirms))
}
func (l *Logger) AverageDeclinedRouteMessageReceives() float64 {
	return float64(l.CountDeclinedRouteMessageReceives()) / float64(len(l.deniedRouteMessageReceives))
}
func (l *Logger) AverageDeclinedRouteMessageConfirms() float64 {
	return float64(l.CountDeclinedRouteMessageConfirms()) / float64(len(l.deniedRouteMessageConfirms))
}

func (l *Logger) AverageDownloadMessages() float64 {
	return float64(l.CountDownloadMessages()) / float64(len(l.downloadMessages))
}

func (l *Logger) AverageDownloadPath() float64 {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	sum := 0
	for id := range l.downloadMessages {
		sum += l.countDownloadPathWithoutLock(id)
	}

	return float64(sum) / float64(len(l.downloadMessages))
}
