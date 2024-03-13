package basenode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	netp "dirs/simulation/pkg/network"
	"dirs/simulation/pkg/utils"
	"sync"
)

type BaseNode struct {
	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex
	// Requests holds unaswered yet request
	requests     []_Request
	requestsLock sync.Mutex

	bandwidthManager *bmp.BandwidthManager
	network          *netp.Network[BaseNode]
}

func (n *BaseNode) Receive(key string, val string) {

	n.PutInStore(key, val)

	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = utils.Filter(n.requests, func(r _Request, i int) bool {
		if r.key == key {
			n.RegisterDownload(key, val, r.from)
			return false
		}
		return true
	})
}

func (n *BaseNode) Ask(key string, from *BaseNode) {

	val, ok := n.GetFromStore(key)

	if ok {
		n.RegisterDownload(key, val, from)
	} else {
		friends := n.network.GetFriends(n)

		if len(friends) < 2 {
			return
		}

		n.AddRequest(key, from)

		for _, friend := range friends {
			if friend == from {
				continue
			}
			go friend.Ask(key, n)
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
