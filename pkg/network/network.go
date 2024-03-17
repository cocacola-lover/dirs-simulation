package network

import (
	fp "dirs/simulation/pkg/fundamentals"
	gp "dirs/simulation/pkg/graph"
	lp "dirs/simulation/pkg/logger"
)

type Tunnel struct {
	Width  int
	Length int
}

type Network struct {
	gp.Graph[Tunnel]
	nodes     []fp.INode
	phoneBook map[fp.INode]int
	Logger    *lp.Logger
}

func (net Network) GetTunnel(node1, node2 fp.INode) Tunnel {
	tunnel, ok := net.HasPath(net.phoneBook[node1], net.phoneBook[node2])

	if ok {
		return tunnel
	} else {
		panic("Tried to get Tunnel between two nodes that do not have a tunnel")
	}
}

func (net Network) StringLogger() string {
	return net.Logger.String(net.phoneBook)
}

func (net Network) GetFriends(node fp.INode) []fp.INode {
	ans := []fp.INode{}

	indsOfFriends, _ := net.GetPaths(net.phoneBook[node])

	for _, i := range indsOfFriends {
		ans = append(ans, net.nodes[i])
	}

	return ans
}

func (net Network) Get(i int) fp.INode {
	return net.nodes[i]
}

func (net Network) String() string {
	str := net.Graph.String()
	return str
}
