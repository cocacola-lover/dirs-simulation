package network

import (
	crp "dirs/simulation/pkg/controlledRandom"
	gp "dirs/simulation/pkg/graph"
	lnp "dirs/simulation/pkg/loggerNode"
	lp "dirs/simulation/pkg/nLogger"
	np "dirs/simulation/pkg/node"
	searchernode "dirs/simulation/pkg/searcherNode"
)

func _NewWithoutGraphNetwork(initNode func(net *Network, i int) *lnp.LoggerNode, size int, logger *lp.Logger) *Network {

	var net Network

	if logger != nil {
		net = Network{
			nodes:     make([]*lnp.LoggerNode, size),
			phoneBook: make(map[np.INode]int),
			Logger:    logger,
		}
	} else {
		net = Network{
			nodes:     make([]*lnp.LoggerNode, size),
			phoneBook: make(map[np.INode]int),
			Logger:    lp.NewLogger(),
		}
	}

	for i := 0; i < size; i++ {
		net.nodes[i] = initNode(&net, i)
		net.phoneBook[net.nodes[i]] = i
	}

	return &net
}

func NewEmptyNetwork(initNode func(net *Network, i int) *lnp.LoggerNode, size int, logger *lp.Logger) *Network {
	network := _NewWithoutGraphNetwork(initNode, size, logger)

	network.Graph = gp.NewGraph[Tunnel](size)

	return network
}

func NewRandomNetwork(initNode func(net *Network, i int) *lnp.LoggerNode, size int, degree int, logger *lp.Logger) *Network {

	network := _NewWithoutGraphNetwork(initNode, size, logger)

	network.Graph = gp.NewRandomConnectedGraph[Tunnel](size, degree, func() Tunnel {
		return Tunnel{Width: crp.Rand.Intn(10) + 1, Length: crp.Rand.Intn(10) + 1}
	})

	return network
}

// Test Network that is populated entirely by base nodes.
func NewBaseNetwork(size, degree int, logger *lp.Logger) *Network {
	return NewRandomNetwork(func(net *Network, i int) *lnp.LoggerNode {
		bn := np.NewNode(
			crp.Rand.Intn(10)+1,
			crp.Rand.Intn(10)+1,
			nil, nil,
		)

		n := lnp.NewLoggerNode(searchernode.NewSearchNode(bn), net.Logger)

		bn.SetOuterFunctions(
			func() []np.INode {
				return net.GetFriends(n)
			}, func(with np.INode) (int, int) {
				return net.GetTunnel(n, with)
			}, nil,
		)

		return n
	}, size, degree, logger)
}

// Test Network that is populated by failing nodes.
func NewFailingNetwork(size, degree int, logger *lp.Logger) *Network {
	return NewRandomNetwork(func(net *Network, i int) *lnp.LoggerNode {
		bn := np.NewNode(
			crp.Rand.Intn(10)+1,
			crp.Rand.Intn(10)+1,
			nil, nil,
		)

		n := lnp.NewLoggerNode(searchernode.NewSearchNode(bn), net.Logger)

		bn.SetOuterFunctions(
			func() []np.INode {
				return net.GetFriends(n)
			}, func(with np.INode) (int, int) {
				return net.GetTunnel(n, with)
			}, func(method np.Method) float64 {
				if method == np.ReceiveDownloadMethod {
					return 0.1
				}
				return 0
			},
		)

		return n
	}, size, degree, logger)
}
