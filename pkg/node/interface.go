package node

import bmp "dirs/simulation/pkg/bandwidthManager"

type INode interface {
	ReceiveRouteMessage(id int, key string, from INode) bool
	ConfirmRouteMessage(id int, from INode) bool
	TimeoutRouteMessage(id int, from INode) bool
	// Returns whether message was accepted or rejected ^

	ReceiveDownloadMessage(id int, key string, from INode)
	ConfirmDownloadMessage(id int, val string, from INode)
	Bm() *bmp.BandwidthManager
	SetSelfAddress(n INode)
}
