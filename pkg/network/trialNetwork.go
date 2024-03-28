package network

import (
	"time"
)

type SearchRequest struct {
	Id  int
	Val string
	Key string

	// Number between 0 and 1
	Popularity        float64
	NumberOfSearchers int
}

type TrialNetwork struct {
	*Network
}

var idCounter int = 0

func (tn *TrialNetwork) GenerateTasks(reqGen func() SearchRequest, intervalGen func() time.Duration, timer time.Duration) {

	start := time.Now()

	for time.Since(start) < timer {
		req := reqGen()

		searchers, havers := devideSearchersAndHavers(len(tn.nodes), req)

		for _, each := range havers {
			tn.Get(each).PutVal(req.Key, req.Val)
		}

		for _, each := range searchers {
			go tn.Get(each).ReceiveRouteMessage(idCounter, req.Key, tn.Get(each))
			idCounter++
		}

		time.Sleep(intervalGen())
	}
}

func (tn *TrialNetwork) String() string {
	return tn.Logger.String()
}

func NewTrialNetwork(net *Network) *TrialNetwork {
	return &TrialNetwork{Network: net}
}
