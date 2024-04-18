package node

import bmp "dirs/simulation/pkg/bandwidthManager"

type INode interface {
	StartSearch(key string) (int, bool)

	ReceiveRouteMessage(id int, key string, from INode) bool
	ConfirmRouteMessage(id int, from INode) bool
	// Returns whether message was accepted or rejected ^

	Fail()
	RetryMessages(ids []Request) []int
	ReceiveFaultMessage(from INode, about []int) []int

	ReceiveDownloadMessage(id int, key string, from INode)
	ConfirmDownloadMessage(id int, val string, from INode)

	Bm() *bmp.BandwidthManager
	HasFailed() bool
	Close()

	SetSelfAddress(n INode)
	GetSelfAddress() INode

	PutKey(key, val string)
	ReceivedKey(key string) (string, bool)

	AddToStore(key, val string)
	HasInStore(key string) (string, bool)
}

type Method uint

const (
	ReceiveRouteMethod Method = iota
	ConfirmRouteMethod
	ReceiveDownloadMethod
	ConfirmDownloadMethod
)
