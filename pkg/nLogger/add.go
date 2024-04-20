package nlogger

import (
	"dirs/simulation/pkg/node"
	"sync"
	"sync/atomic"
	"time"
)

func (l *Logger) AddFailedNode() {
	atomic.StoreInt32(&l.failedNode, atomic.LoadInt32(&l.failedNode)+1)
}

func addToMapMapArrWithLock(id int, from node.INode, to node.INode, store map[int]int, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()

	v, ok := store[id]
	if !ok {
		store[id] = 1
	} else {
		store[id] = v + 1
	}
}

func (l *Logger) AddRouteMessageReceive(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.routeMessageReceives, &l.rmrLock)
}

func (l *Logger) AddFaultMessageReceive() {
	atomic.StoreInt32(&l.faultMessageReceives, atomic.LoadInt32(&l.faultMessageReceives)+1)
}

func (l *Logger) AddRouteMessageConfirm(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.routeMessageConfirms, &l.rmcLock)
}
func (l *Logger) AddDeniedRouteMessageReceive(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.deniedRouteMessageReceives, &l.rmrdLock)
}
func (l *Logger) AddDeniedRouteMessageConfirm(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.deniedRouteMessageConfirms, &l.rmcdLock)
}

func (l *Logger) AddDownloadMessage(id int, from node.INode, to node.INode) {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	_, ok := l.downloadMessages[id]
	if !ok {
		l.downloadMessages[id] = make(map[node.INode]node.INode)
	}

	l.downloadMessages[id][to] = from
}

// Should be used synchronously
func (l *Logger) StartSearch(id int, n node.INode) {
	timestamp := make([]time.Time, 2)
	timestamp[0] = time.Now()

	go func() {
		l.setLock.Lock()
		l.seTimestamps[id] = timestamp
		l.setLock.Unlock()

		l.dLock.Lock()
		l.startedSearches[id] = n
		l.dLock.Unlock()
	}()
}

// Should be used synchronously
func (l *Logger) EndSearch(id int) {
	endTimestamp := time.Now()

	go func() {
		l.setLock.Lock()
		l.seTimestamps[id][1] = endTimestamp
		l.setLock.Unlock()
	}()
}
