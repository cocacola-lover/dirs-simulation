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
	storeLock sync.Mutex
	// Requests holds unaswered yet request
	requests     []_Request
	requestsLock sync.Mutex

	bandwidthManager *bmp.BandwidthManager
	network          *netp.Network[BaseNode]
}

func (n *BaseNode) RegisterDownload(key string, val string, with *BaseNode) {
	n.bandwidthManager.RegisterDownload(
		len(val),
		with.bandwidthManager,
		n.network.GetTunnelWidth(n, with),
		func() {
			with.Receive(key, val)
		},
	)
}

func (n *BaseNode) Receive(key string, val string) {

	utils.WithLockedNoResult(&n.storeLock, func() {
		n.store[key] = val
	})

	utils.WithLockedNoResult(&n.requestsLock, func() {
		n.requests = utils.Filter(n.requests, func(r _Request, i int) bool {
			if r.key == key {
				n.RegisterDownload(key, val, r.from)
				return false
			}
			return true
		})
	})
}

func (n *BaseNode) Ask(key string, from *BaseNode) {

	n.storeLock.Lock()
	val, ok := n.store[key]
	n.storeLock.Unlock()

	if ok {
		n.RegisterDownload(key, val, from)
	} else {

		friends := n.network.GetFriends(n)

		if len(friends) == 0 {
			return
		}

		if len(friends) == 1 && friends[0] == from {
			return
		}

		utils.WithLockedNoResult(&n.requestsLock, func() {
			n.requests = append(n.requests, _Request{
				from: from, key: key,
			})
		})

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
