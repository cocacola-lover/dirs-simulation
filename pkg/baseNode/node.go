package basenode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	"dirs/simulation/pkg/utils"
	"sync"
)

type BaseNode struct {
	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex
	// Requests holds unaswered yet request
	requests     []IMessage
	requestsLock sync.RWMutex

	bandwidthManager *bmp.BandwidthManager

	// OuterGetterFunctions

	getFriends func() []*BaseNode
	getTunnel  func(with *BaseNode) (int, int)

	// Watchers

	watchPutInStore       func(m IMessage, val string)
	watchRegisterDownload func(m IMessage, val string)
}

func (n *BaseNode) Receive(newm IMessage, val string) {

	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = utils.Filter(n.requests, func(m IMessage, i int) bool {
		if m.Key() == newm.Key() {
			if m.From() == n {
				n.PutInStore(m, val)
			} else {
				n.RegisterDownload(m, val)
			}

			return false
		}
		return true
	})
}

func (n *BaseNode) Ask(m IMessage) {

	if n.HasMessage(m) || !m.IsValid() {
		return
	}

	val, ok := n.GetFromStore(m.Key())

	if ok {
		if m.From() != n {
			n.RegisterDownload(m, val)
		}
	} else {
		toAsk := n.WhoToAsk(m)
		if len(toAsk) == 0 {
			return
		}

		n.AddRequest(m)

		for _, friend := range toAsk {
			go friend.Ask(Resend(m, n))
		}
	}
}
