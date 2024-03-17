package basenode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	fp "dirs/simulation/pkg/fundamentals"
)

// Initialization functionS
func NewBaseNode(maxDownload int, maxUpload int) *BaseNode {
	return &BaseNode{
		bandwidthManager: bmp.NewBandwidthManager(maxDownload, maxUpload),
		store:            make(map[string]string),
	}
}

// getTunnel returns (tunnelWidth, tunnelLength)
func (n *BaseNode) SetGetters(getFriends func() []fp.INode, getTunnel func(with fp.INode) (int, int)) *BaseNode {
	n.getFriends = getFriends
	n.getTunnel = getTunnel

	return n
}

func (n *BaseNode) SetWatchers(watchPutInStore func(m fp.IMessage, me fp.INode), watchRegisterDownload func(m fp.IMessage, me fp.INode)) *BaseNode {
	n.watchPutInStore = watchPutInStore
	n.watchRegisterDownload = watchRegisterDownload

	return n
}

func (n *BaseNode) InitStore(store map[string]string) {
	n.storeLock.Lock()
	defer n.storeLock.Unlock()

	for k, v := range store {
		n.store[k] = v
	}
}
