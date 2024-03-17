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
	requests     []_Request
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

	n.requests = utils.Filter(n.requests, func(m _Request, i int) bool {
		if m.Key() == newm.Key() {
			if m.From() == n {
				n.registerInStore(m, val)
			} else if m.From().IsInterestedIn(m.Key()) {
				n.registerDownload(m, val)
			}
			m.Done(n)

			return false
		}
		return true
	})
}

func (n *BaseNode) Ask(m fp.IMessage) {

	if n.hasMessage(m) {
		return
	}

	val, ok := n.getFromStore(m.Key())

	if ok {
		if m.From() != n {
			m.Done(n)
			n.registerDownload(m, val)
		}
	} else {
		toAsk := n.whoToAsk(m)
		if len(toAsk) == 0 {
			return
		}

		n.addRequest(_NewRequest(m, toAsk))

		for _, friend := range toAsk {
			go friend.Ask(m.Resend(n))
		}
	}
}

func (n *BaseNode) IsInterestedIn(key string) bool {
	n.requestsLock.RLock()
	defer n.requestsLock.RUnlock()

	for _, r := range n.requests {
		if r.Key() == key {
			return true
		}
	}

	return false
}

func (n *BaseNode) StopSearch(id int, from fp.INode) {

	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	for i := 0; i < len(n.requests); i++ {
		if n.requests[i].Id() == id && n.requests[i].From() == from {
			go n.requests[i].stopSearch(n)
			n.requests[i] = n.requests[len(n.requests)-1]
			n.requests = n.requests[:len(n.requests)-1]
			break
		}
	}
}
