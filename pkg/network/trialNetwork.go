package network

import (
	"sync"
)

type SearchRequest struct {
	Val string
	Key string

	// Number between 0 and 1
	Popularity        float64
	NumberOfSearchers int
}

type TrialNetwork struct {
	*Network
}

func (tn *TrialNetwork) WaitToFinishAllSearchers() chan bool {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(tn.nodes))

	for _, each := range tn.nodes {
		go each.WaitToFinishAllSearches(waitGroup)
	}

	exitCh := make(chan bool)

	go func() {
		waitGroup.Wait()
		close(exitCh)
	}()

	return exitCh
}

// Returns havers and searchers
func (tn *TrialNetwork) RunRequest(req SearchRequest) ([]int, []int) {

	searchers, havers := devideSearchersAndHavers(len(tn.nodes), req)

	// fmt.Printf("Searchers are %v\n", searchers)
	// fmt.Printf("Havers are %v\n", havers)

	for _, each := range havers {
		tn.Get(each).INode.AddToStore(req.Key, req.Val)
	}

	for _, each := range searchers {
		go tn.Get(each).StartSearchAndWatch(req.Key)
	}

	return havers, searchers

}

func (tn *TrialNetwork) String() string {
	return tn.Logger.String()
}

func (tn *TrialNetwork) Close() {
	for _, n := range tn.nodes {
		n.Close()
	}
}

func NewTrialNetwork(net *Network) *TrialNetwork {
	return &TrialNetwork{Network: net}
}
