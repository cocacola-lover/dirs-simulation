package nlogger

import (
	"dirs/simulation/pkg/node"
	"sync"
)

func addToMapMapArrWithLock(id int, from node.INode, to node.INode, store map[int]map[node.INode][]node.INode, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()

	dict, ok := store[id]
	if !ok {
		store[id] = make(map[node.INode][]node.INode)
		dict = store[id]
	}

	dict[from] = append(dict[from], to)
}

func (l *Logger) AddRouteMessageReceive(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.routeMessageReceives, &l.rmrLock)
}

func (l *Logger) AddRouteMessageTimeout(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.routeMessageTimeouts, &l.rmtLock)
}

func (l *Logger) AddRouteMessageConfirm(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.routeMessageConfirms, &l.rmcLock)
}
func (l *Logger) AddDeniedRouteMessageReceive(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.deniedRouteMessageReceives, &l.rmrdLock)
}
func (l *Logger) AddDeniedRouteMessageTimeout(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.deniedRouteMessageTimeouts, &l.rmtdLock)
}
func (l *Logger) AddDeniedRouteMessageConfirm(id int, from node.INode, to node.INode) {
	addToMapMapArrWithLock(id, from, to, l.deniedRouteMessageConfirms, &l.rmcdLock)
}

func (l *Logger) AddDownloadMessage(id int, from node.INode) {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	l.downloadMessages[id] = append(l.downloadMessages[id], from)
}
