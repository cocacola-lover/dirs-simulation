package fundamentals

import bmp "dirs/simulation/pkg/bandwidthManager"

type INode interface {
	BandwidthManager() *bmp.BandwidthManager
	Receive(m IMessage, val string)
	Ask(m IMessage)
	IsInterestedIn(key string) bool
	StopSearch(id int, from INode)

	InitStore(store map[string]string)
}
