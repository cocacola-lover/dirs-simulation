package network

import (
	crp "dirs/simulation/pkg/controlledRandom"
	gp "dirs/simulation/pkg/graph"
	lnp "dirs/simulation/pkg/loggerNode"
	lp "dirs/simulation/pkg/nLogger"
	np "dirs/simulation/pkg/node"
)

func _NewWithoutGraphNetwork(initNode func(net *Network, i int) np.INode, size int) *Network {
	net := Network{
		nodes:     make([]np.INode, size),
		phoneBook: make(map[np.INode]int),
		Logger:    lp.NewLogger(),
	}

	for i := 0; i < size; i++ {
		net.nodes[i] = initNode(&net, i)
		net.phoneBook[net.nodes[i]] = i
	}

	return &net
}

func NewEmptyNetwork(initNode func(net *Network, i int) np.INode, size int) *Network {
	network := _NewWithoutGraphNetwork(initNode, size)

	network.Graph = gp.NewGraph[Tunnel](size)

	return network
}

func NewRandomNetwork(initNode func(net *Network, i int) np.INode, size int, degree int) *Network {

	network := _NewWithoutGraphNetwork(initNode, size)

	network.Graph = gp.NewRandomConnectedGraph[Tunnel](size, degree, func() Tunnel {
		return Tunnel{Width: crp.Rand.Intn(10) + 1, Length: crp.Rand.Intn(10) + 1}
	})

	return network
}

// Test Network that is populated entirely of base nodes.
func NewTestSearchNetwork(size int, degree int, requests ...SearchRequest) (*Network, [][]int, [][]int) {

	var perRequestSearchers [][]int
	var perRequestHavers [][]int

	for _, r := range requests {
		searchers, havers := devideSearchersAndHavers(size, r)
		perRequestSearchers = append(perRequestSearchers, searchers)
		perRequestHavers = append(perRequestHavers, havers)
	}

	perRequestHaversMaps := []map[int]bool{}
	for i := range perRequestHavers {
		perRequestHaversMaps = append(perRequestHaversMaps, make(map[int]bool))
		for _, v := range perRequestHavers[i] {
			perRequestHaversMaps[i][v] = true
		}
	}

	network := NewRandomNetwork(func(net *Network, i int) np.INode {
		bn := np.NewNode(
			crp.Rand.Intn(10)+1,
			crp.Rand.Intn(10)+1,
			nil, nil,
		)

		n := lnp.NewLoggerNode(bn, net.Logger)

		bn.SetOuterFunctions(
			func() []np.INode {
				return net.GetFriends(n)
			}, func(with np.INode) (int, int) {
				return net.GetTunnel(n, with)
			},
		)

		for j, havers := range perRequestHaversMaps {
			if _, ok := havers[i]; ok {
				bn.PutVal(requests[j].Key, requests[j].Val)
			}
		}

		return n
	}, size, degree)

	return network, perRequestSearchers, perRequestHavers
}

// Private -----------------------------------------------------------------------

type SearchRequest struct {
	Val string
	Key string

	// Number between 0 and 1
	Popularity        float64
	NumberOfSearchers int
}

func devideSearchersAndHavers(size int, request SearchRequest) ([]int, []int) {
	hasInStore := []int{}

	for i := 0; i < size; i++ {
		if crp.Rand.Float64() <= request.Popularity {
			hasInStore = append(hasInStore, i)
		}
	}

	// But at least 1 has info
	if len(hasInStore) == 0 {
		hasInStore = append(hasInStore, crp.Rand.Intn(size))
	}

	searchers := []int{}
	for i := 0; i < request.NumberOfSearchers; i++ {
		ind := crp.Rand.Intn(size - len(hasInStore) - len(searchers))

		jumpOver := 0
		for h, s := 0, 0; h+s < len(hasInStore)+len(searchers); {
			var v int
			if len(searchers) > s && hasInStore[h] > searchers[s] {
				v = searchers[s]
				s++
			} else {
				v = hasInStore[h]
				h++
			}

			if v > ind+jumpOver {
				break
			}
			jumpOver++
		}

		searchers = append(searchers, ind+jumpOver)
	}

	return searchers, hasInStore
}
