package nlogger

import (
	"dirs/simulation/pkg/node"
	"sync"
)

type Logger struct {
	// Send in search by id from higher node to lower nodes
	routeMessageReceives map[int]map[node.INode][]node.INode
	// Send in search by id from higher node to lower nodes
	routeMessageTimeouts map[int]map[node.INode][]node.INode
	// Send in search by id from lower node to higher nodes
	routeMessageConfirms map[int]map[node.INode][]node.INode

	downloadMessages map[int][]node.INode

	rmrLock sync.Mutex
	rmtLock sync.Mutex
	rmcLock sync.Mutex
	dLock   sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{
		routeMessageReceives: make(map[int]map[node.INode][]node.INode),
		routeMessageTimeouts: make(map[int]map[node.INode][]node.INode),
		routeMessageConfirms: make(map[int]map[node.INode][]node.INode),
		downloadMessages:     make(map[int][]node.INode),
	}
}
