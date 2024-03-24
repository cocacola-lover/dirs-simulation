package nlogger

import (
	np "dirs/simulation/pkg/node"
	"sync"
)

type Logger struct {
	// Send in search by id from higher node to lower nodes
	routeMessageReceives       map[int]map[np.INode][]np.INode
	deniedRouteMessageReceives map[int]map[np.INode][]np.INode
	// Send in search by id from higher node to lower nodes
	routeMessageTimeouts       map[int]map[np.INode][]np.INode
	deniedRouteMessageTimeouts map[int]map[np.INode][]np.INode
	// Send in search by id from lower node to higher nodes
	routeMessageConfirms       map[int]map[np.INode][]np.INode
	deniedRouteMessageConfirms map[int]map[np.INode][]np.INode

	downloadMessages map[int][]np.INode

	rmrLock  sync.Mutex
	rmtLock  sync.Mutex
	rmcLock  sync.Mutex
	rmrdLock sync.Mutex
	rmtdLock sync.Mutex
	rmcdLock sync.Mutex
	dLock    sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{
		routeMessageReceives:       make(map[int]map[np.INode][]np.INode),
		routeMessageTimeouts:       make(map[int]map[np.INode][]np.INode),
		routeMessageConfirms:       make(map[int]map[np.INode][]np.INode),
		deniedRouteMessageReceives: make(map[int]map[np.INode][]np.INode),
		deniedRouteMessageTimeouts: make(map[int]map[np.INode][]np.INode),
		deniedRouteMessageConfirms: make(map[int]map[np.INode][]np.INode),
		downloadMessages:           make(map[int][]np.INode),
	}
}
