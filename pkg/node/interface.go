package node

import bmp "dirs/simulation/pkg/bandwidthManager"

type INode interface {
	StartSearch(key string) int

	ReceiveRouteMessage(id int, key string, from INode) bool
	ConfirmRouteMessage(id int, from INode) bool
	// Returns whether message was accepted or rejected ^

	Fail()
	RetryMessages(ids []Request) []int
	ReceiveFaultMessage(from INode, about []int) []int

	ReceiveDownloadMessage(id int, key string, from INode)
	ConfirmDownloadMessage(id int, val string, from INode)

	Bm() *bmp.BandwidthManager
	Close()

	SetSelfAddress(n INode)
	GetSelfAddress() INode

	PutVal(key, val string)
	HasKey(key string) (string, bool)
}

type Method uint

const (
	ReceiveRouteMethod Method = iota
	ConfirmRouteMethod
	ReceiveDownloadMethod
	ConfirmDownloadMethod
)
