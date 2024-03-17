package network

import (
	bnp "dirs/simulation/pkg/baseNode"
	crp "dirs/simulation/pkg/controlledRandom"
	fp "dirs/simulation/pkg/fundamentals"
)

func initBaseNode(net *Network, i int, maxUpload int, maxDownload int) *bnp.BaseNode {

	logMessage := func(m fp.IMessage, me fp.INode) {
		net.Logger.AddMessage(m, me)
	}

	node := bnp.NewBaseNode(maxDownload, maxUpload)
	node = node.SetGetters(func() []fp.INode {
		return net.GetFriends(node)
	}, func(with fp.INode) (int, int) {
		tunnel := net.GetTunnel(node, with)
		return tunnel.Width, tunnel.Length
	}).SetWatchers(logMessage, logMessage)

	return node
}

func FactoryInitBaseNode() func(net *Network, i int) fp.INode {
	return func(net *Network, i int) fp.INode {
		return initBaseNode(net, i, 1, 1)
	}
}

// Min-max corresponds to borders of random maxUploads and maxDownloads
func FactoryInitRandomBaseNode(min, max int) func(net *Network, i int) fp.INode {
	return func(net *Network, i int) fp.INode {
		return initBaseNode(net, i, crp.Rand.Intn(max)+min, crp.Rand.Intn(max)+min)
	}
}
