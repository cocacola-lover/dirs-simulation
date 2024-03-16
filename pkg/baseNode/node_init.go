package basenode

import bmp "dirs/simulation/pkg/bandwidthManager"

// Initialization functionS
func NewBaseNode(maxDownload int, maxUpload int) *BaseNode {
	return &BaseNode{
		bandwidthManager: bmp.NewBandwidthManager(maxDownload, maxUpload),
		store:            make(map[string]string),
	}
}

// getTunnel returns (tunnelWidth, tunnelLength)
func (n *BaseNode) SetGetters(getFriends func() []*BaseNode, getTunnel func(with *BaseNode) (int, int)) *BaseNode {
	n.getFriends = getFriends
	n.getTunnel = getTunnel

	return n
}

func (n *BaseNode) SetWatchers(watchPutInStore func(m IMessage, val string), watchRegisterDownload func(m IMessage, val string)) *BaseNode {
	n.watchPutInStore = watchPutInStore
	n.watchRegisterDownload = watchRegisterDownload

	return n
}
