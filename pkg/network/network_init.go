package network

import (
	crp "dirs/simulation/pkg/controlledRandom"
	gp "dirs/simulation/pkg/graph"
	lazynode "dirs/simulation/pkg/lazyNode"
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

// Test Network that is populated entirely by base nodes.
func NewBaseNetwork(size, degree int) *Network {
	return NewRandomNetwork(func(net *Network, i int) np.INode {
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

		return n
	}, size, degree)
}

// Test Network that is populated by base nodes and lazy nodes.
func NewLazyNetwork(size int, degree int, lazyPop float64) *Network {
	return NewRandomNetwork(func(net *Network, i int) np.INode {
		var n np.INode

		if lazyPop >= crp.Rand.Float64() {
			bn := lazynode.NewLazyNode(
				crp.Rand.Intn(10)+1,
				crp.Rand.Intn(10)+1,
				nil, nil,
			)

			n = lnp.NewLoggerNode(bn, net.Logger)

			bn.SetOuterFunctions(
				func() []np.INode {
					return net.GetFriends(n)
				}, func(with np.INode) (int, int) {
					return net.GetTunnel(n, with)
				},
			)
		} else {
			bn := np.NewNode(
				crp.Rand.Intn(10)+1,
				crp.Rand.Intn(10)+1,
				nil, nil,
			)

			n = lnp.NewLoggerNode(bn, net.Logger)

			bn.SetOuterFunctions(
				func() []np.INode {
					return net.GetFriends(n)
				}, func(with np.INode) (int, int) {
					return net.GetTunnel(n, with)
				},
			)
		}

		return n
	}, size, degree)
}
