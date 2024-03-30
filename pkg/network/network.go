package network

import (
	gp "dirs/simulation/pkg/graph"
	lnp "dirs/simulation/pkg/loggerNode"
	nlogger "dirs/simulation/pkg/nLogger"
	np "dirs/simulation/pkg/node"
)

type Tunnel struct {
	Width  int
	Length int
}

type Network struct {
	gp.Graph[Tunnel]
	nodes     []*lnp.LoggerNode
	phoneBook map[np.INode]int
	Logger    *nlogger.Logger
}

// Returns tunnel.Width and tunnel.Length
func (net Network) GetTunnel(node1, node2 np.INode) (int, int) {
	tunnel, ok := net.HasPath(net.phoneBook[node1], net.phoneBook[node2])

	if ok {
		return tunnel.Width, tunnel.Length
	} else {
		panic("Tried to get Tunnel between two nodes that do not have a tunnel")
	}
}

func (net Network) GetFriends(node np.INode) []np.INode {
	ans := []np.INode{}

	indsOfFriends, _ := net.GetPaths(net.phoneBook[node])

	for _, i := range indsOfFriends {
		ans = append(ans, net.nodes[i])
	}

	return ans
}

func (net Network) Get(i int) *lnp.LoggerNode {
	return net.nodes[i]
}

func (net Network) String() string {
	str := net.Graph.String()
	return str
}

func (net Network) LoggerStringById(id int, lead *lnp.LoggerNode) string {
	return net.Logger.StringByIdVerbose(id, lead, net.phoneBook)
}
