package nlogger

import "dirs/simulation/pkg/node"

func (l *Logger) AddRouteMessageReceive(id int, from node.INode, to node.INode) {
	l.rmrLock.Lock()
	defer l.rmrLock.Unlock()

	dict, ok := l.routeMessageReceives[id]
	if !ok {
		l.routeMessageReceives[id] = make(map[node.INode][]node.INode)
		dict = l.routeMessageReceives[id]
	}

	dict[from] = append(dict[from], to)
}

func (l *Logger) AddRouteMessageTimeout(id int, from node.INode, to node.INode) {
	l.rmtLock.Lock()
	defer l.rmtLock.Unlock()

	dict, ok := l.routeMessageTimeouts[id]
	if !ok {
		l.routeMessageTimeouts[id] = make(map[node.INode][]node.INode)
		dict = l.routeMessageTimeouts[id]
	}

	dict[from] = append(dict[from], to)
}

func (l *Logger) AddRouteMessageConfirms(id int, from node.INode, to node.INode) {
	l.rmcLock.Lock()
	defer l.rmcLock.Unlock()

	dict, ok := l.routeMessageConfirms[id]
	if !ok {
		l.routeMessageConfirms[id] = make(map[node.INode][]node.INode)
		dict = l.routeMessageConfirms[id]
	}

	dict[from] = append(dict[from], to)
}

func (l *Logger) AddDownloadMessage(id int, from node.INode) {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	l.downloadMessages[id] = append(l.downloadMessages[id], from)
}
