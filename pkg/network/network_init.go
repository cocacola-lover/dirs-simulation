package network

import (
	crp "dirs/simulation/pkg/controlledRandom"
	fp "dirs/simulation/pkg/fundamentals"
	gp "dirs/simulation/pkg/graph"
	lp "dirs/simulation/pkg/logger"
)

func _NewWithoutGraphNetwork(initNode func(net *Network, i int) fp.INode, size int) *Network {
	net := Network{
		nodes:     make([]fp.INode, size),
		phoneBook: make(map[fp.INode]int),
		Logger:    lp.NewLogger(),
	}

	for i := 0; i < size; i++ {
		net.nodes[i] = initNode(&net, i)
		net.phoneBook[net.nodes[i]] = i
	}

	return &net
}

func NewEmptyNetwork(initNode func(net *Network, i int) fp.INode, size int) *Network {
	network := _NewWithoutGraphNetwork(initNode, size)

	network.Graph = gp.NewGraph[Tunnel](size)

	return network
}

func NewRandomNetwork(initNode func(net *Network, i int) fp.INode, size int, degree int) *Network {

	network := _NewWithoutGraphNetwork(initNode, size)

	network.Graph = gp.NewRandomConnectedGraph[Tunnel](size, degree, func() Tunnel {
		return Tunnel{Width: crp.Rand.Intn(10) + 1, Length: crp.Rand.Intn(10) + 1}
	})

	return network
}
