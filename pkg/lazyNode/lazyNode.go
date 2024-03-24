package lazynode

import "dirs/simulation/pkg/node"

type LazyNode struct {
	*node.Node
}

func (ln *LazyNode) ReceiveRouteMessage(id int, key string, from node.INode) bool {
	if ln.GetSelfAddress() != from {
		return false
	} else {
		return ln.Node.ReceiveRouteMessage(id, key, from)
	}
}

func NewLazyNode(maxDownload int, maxUpload int, getNetworkFriends func() []node.INode, getNetworkTunnel func(with node.INode) (int, int)) *LazyNode {
	n := &LazyNode{Node: node.NewNode(maxDownload, maxUpload, getNetworkFriends, getNetworkTunnel)}
	n.SetSelfAddress(n)
	return n
}

func (n *LazyNode) SetOuterFunctions(getNetworkFriends func() []node.INode, getNetworkTunnel func(with node.INode) (int, int)) *LazyNode {
	n.Node.SetOuterFunctions(getNetworkFriends, getNetworkTunnel)
	return n
}
