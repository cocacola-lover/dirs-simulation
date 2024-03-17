package fundamentals

import bmp "dirs/simulation/pkg/bandwidthManager"

type INode interface {
	BandwidthManager() *bmp.BandwidthManager
	Receive(m IMessage, val string)
	Ask(m IMessage)
	IsInterestedIn(key string) bool

	InitStore(store map[string]string)
}
