package nlogger

import (
	np "dirs/simulation/pkg/node"
	"sync"
	"time"
)

type Logger struct {
	// Start-End timestamps
	seTimestamps map[int][]time.Time
	// Send in search by id from higher node to lower nodes
	routeMessageReceives       map[int]map[np.INode][]np.INode
	deniedRouteMessageReceives map[int]map[np.INode][]np.INode
	// Send in search by id from higher node to lower nodes
	faultMessageReceives map[int]map[np.INode][]np.INode
	// Send in search by id from lower node to higher nodes
	routeMessageConfirms       map[int]map[np.INode][]np.INode
	deniedRouteMessageConfirms map[int]map[np.INode][]np.INode

	// Started searches
	startedSearches map[int]np.INode
	// Tracks To - key, from - value
	downloadMessages map[int]map[np.INode]np.INode

	setLock sync.Mutex
	rmrLock sync.Mutex
	fmrLock sync.Mutex
	rmcLock sync.Mutex

	rmrdLock sync.Mutex
	rmcdLock sync.Mutex

	// Locks both startedSearch and downloadMessages
	dLock sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{
		routeMessageReceives:       make(map[int]map[np.INode][]np.INode),
		faultMessageReceives:       make(map[int]map[np.INode][]np.INode),
		routeMessageConfirms:       make(map[int]map[np.INode][]np.INode),
		deniedRouteMessageReceives: make(map[int]map[np.INode][]np.INode),
		deniedRouteMessageConfirms: make(map[int]map[np.INode][]np.INode),
		downloadMessages:           make(map[int]map[np.INode]np.INode),
		seTimestamps:               make(map[int][]time.Time),
		startedSearches:            make(map[int]np.INode),
	}
}
