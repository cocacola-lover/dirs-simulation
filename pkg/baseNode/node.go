package basenode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	fp "dirs/simulation/pkg/fundamentals"
	"dirs/simulation/pkg/utils"
	"sync"
)

type BaseNode struct {
	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex
	// Requests holds unaswered yet request
	requests     []fp.IMessage
	requestsLock sync.RWMutex

	bandwidthManager *bmp.BandwidthManager

	// OuterGetterFunctions

	getFriends func() []fp.INode
	getTunnel  func(with fp.INode) (int, int)

	// Watchers

	watchPutInStore       func(m fp.IMessage, me fp.INode)
	watchRegisterDownload func(m fp.IMessage, me fp.INode)
}

func (n *BaseNode) Receive(newm fp.IMessage, val string) {

	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = utils.Filter(n.requests, func(m fp.IMessage, i int) bool {
		if m.Key() == newm.Key() {
			if m.From() == n {
				n.registerInStore(m, val)
			} else {
				n.registerDownload(m, val)
			}

			return false
		}
		return true
	})
}

func (n *BaseNode) Ask(m fp.IMessage) {

	if n.hasMessage(m) || !m.IsValid() {
		return
	}

	val, ok := n.getFromStore(m.Key())

	if ok {
		if m.From() != n {
			n.registerDownload(m, val)
		}
	} else {
		toAsk := n.whoToAsk(m)
		if len(toAsk) == 0 {
			return
		}

		n.addRequest(m)

		for _, friend := range toAsk {
			go friend.Ask(m.Resend(n))
		}
	}
}
