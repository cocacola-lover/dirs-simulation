package node

import bmp "dirs/simulation/pkg/bandwidthManager"

type INode interface {
	ReceiveRouteMessage(id int, key string, from INode)
	ConfirmRouteMessage(id int, from INode)
	TimeoutRouteMessage(id int, from INode)
	ReceiveDownloadMessage(id int, key string, from INode)
	ConfirmDownloadMessage(id int, val string, from INode)
	Bm() *bmp.BandwidthManager
}
