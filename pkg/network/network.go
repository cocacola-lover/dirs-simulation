package network

import gp "dirs/simulation/pkg/graph"

type INode interface{}
type Tunnel struct {
	Width  int
	Length int
}

type Network[Node INode] struct {
	gp.Graph[Tunnel]
	nodes     []*Node
	phoneBook map[*Node]int
}

func (net Network[Node]) GetTunnel(node1, node2 *Node) Tunnel {
	tunnel, ok := net.HasPath(net.phoneBook[node1], net.phoneBook[node2])

	if ok {
		return tunnel
	} else {
		panic("Tried to get Tunnel between two nodes that do not have a tunnel")
	}
}

func (net Network[Node]) GetFriends(node *Node) []*Node {
	ans := []*Node{}

	indsOfFriends, _ := net.GetPaths(net.phoneBook[node])

	for _, i := range indsOfFriends {
		ans = append(ans, net.nodes[i])
	}

	return ans
}

func (net Network[Node]) Get(i int) *Node {
	return net.nodes[i]
}

func _NewWithoutGraphNetwork[T INode](initNode func(net *Network[T], i int) *T, size int) *Network[T] {
	net := Network[T]{
		nodes:     make([]*T, size),
		phoneBook: make(map[*T]int),
	}

	for i := 0; i < size; i++ {
		net.nodes[i] = initNode(&net, i)
		net.phoneBook[net.nodes[i]] = i
	}

	return &net
}

func NewEmptyNetwork[T INode](initNode func(net *Network[T], i int) *T, size int) *Network[T] {
	network := _NewWithoutGraphNetwork[T](initNode, size)

	network.Graph = gp.NewGraph[Tunnel](size)

	return network
}

func NewRandomNetwork[T INode](initNode func(net *Network[T], i int) *T, size int, degree int) *Network[T] {

	network := _NewWithoutGraphNetwork[T](initNode, size)

	network.Graph = gp.NewRandomConnectedGraph[Tunnel](size, degree, func() Tunnel {
		return Tunnel{Width: 1, Length: 1}
	})

	return network
}
