package network

import gp "dirs/simulation/pkg/graph"

type INode interface{}

type Graph = gp.Graph

type Network[Node INode] struct {
	Graph
	nodes     []*Node
	phoneBook map[*Node]int
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

	network.Graph = gp.NewGraph(size)

	return network
}

func NewRandomNetwork[T INode](initNode func(net *Network[T], i int) *T, size int, degree int) *Network[T] {

	network := _NewWithoutGraphNetwork[T](initNode, size)

	network.Graph = gp.NewRandomConnectedGraph(size, degree)

	return network
}
