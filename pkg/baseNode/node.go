package basenode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	mp "dirs/simulation/pkg/message"
	netp "dirs/simulation/pkg/network"
	"dirs/simulation/pkg/utils"
	"sync"
)

type BaseNode struct {
	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex
	// Requests holds unaswered yet request
	requests     []mp.BaseMessage[BaseNode]
	requestsLock sync.RWMutex

	bandwidthManager *bmp.BandwidthManager
	network          *netp.Network[BaseNode]
}

func (n *BaseNode) Receive(newm mp.BaseMessage[BaseNode], val string) {

	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = utils.Filter(n.requests, func(m mp.BaseMessage[BaseNode], i int) bool {
		if m.Key == newm.Key {
			if m.From == n {
				n.PutInStore(newm.Key, val)
			} else {
				n.RegisterDownload(m, val)
			}

			go n.network.Logger.AddMessage(m, n)
			return false
		}
		return true
	})
}

func (n *BaseNode) Ask(m mp.BaseMessage[BaseNode]) {

	if n.HasMessage(m) || !m.IsValid() {
		return
	}

	val, ok := n.GetFromStore(m.Key)

	if ok {
		if m.From != n {
			n.network.Logger.AddMessage(m, n)
			n.RegisterDownload(m, val)
		}
	} else {
		friends := n.network.GetFriends(n)

		if len(friends) == 0 || (len(friends) == 1 && friends[0] == m.From) {
			return
		}

		n.AddRequest(m)

		for _, friend := range friends {
			if friend == m.From {
				continue
			}
			go friend.Ask(m.Resend(n))
		}

	}
}

func NewBaseNode(net *netp.Network[BaseNode], maxDownload int, maxUpload int) *BaseNode {
	return &BaseNode{
		bandwidthManager: bmp.NewBandwidthManager(maxDownload, maxUpload),
		store:            make(map[string]string),
		network:          net,
	}
}
