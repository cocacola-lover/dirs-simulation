package nlogger

import (
	"dirs/simulation/pkg/node"
	"sync"
)

func countMapArr(store map[node.INode][]node.INode) int {
	ans := 0

	for _, arr := range store {
		ans += len(arr)
	}

	return ans
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
