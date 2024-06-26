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
	routeMessageReceives       map[int]int
	deniedRouteMessageReceives map[int]int
	// Send in search by id from lower node to higher nodes
	routeMessageConfirms       map[int]int
	deniedRouteMessageConfirms map[int]int

	// Started searches
	startedSearches map[int]np.INode
	// Tracks To - key, from - value
	downloadMessages map[int]map[np.INode]np.INode

	failedNode           int32
	faultMessageReceives int32

	setLock sync.Mutex
	rmrLock sync.Mutex
	rmcLock sync.Mutex

	rmrdLock sync.Mutex
	rmcdLock sync.Mutex

	// Locks both startedSearch and downloadMessages
	dLock sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{
		routeMessageReceives:       make(map[int]int),
		routeMessageConfirms:       make(map[int]int),
		deniedRouteMessageReceives: make(map[int]int),
		deniedRouteMessageConfirms: make(map[int]int),
		downloadMessages:           make(map[int]map[np.INode]np.INode),
		seTimestamps:               make(map[int][]time.Time),
		startedSearches:            make(map[int]np.INode),
	}
}
