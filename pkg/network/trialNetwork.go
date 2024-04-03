package network

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

func (tn *TrialNetwork) WaitToFinishAllSearchers() {
	for _, each := range tn.nodes {
		each.WaitToFinishAllSearches()
	}
}

func (tn *TrialNetwork) RunRequests(reqs []SearchRequest) {
	for _, req := range reqs {
		searchers, havers := devideSearchersAndHavers(len(tn.nodes), req)

		for _, each := range havers {
			tn.Get(each).INode.PutVal(req.Key, req.Val)
		}

		for _, each := range searchers {
			go tn.Get(each).StartSearch(idCounter, req.Key)
			idCounter++
		}
	}
}

func (tn *TrialNetwork) String() string {
	return tn.Logger.String()
}

func NewTrialNetwork(net *Network) *TrialNetwork {
	return &TrialNetwork{Network: net}
}
