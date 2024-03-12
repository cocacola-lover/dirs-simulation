package basenode

import (
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

	network *netp.Network[BaseNode]
}

func (n *BaseNode) Receive(key string, val string) {

	utils.WithLockedNoResult(&n.storeLock, func() {
		n.store[key] = val
	})

	utils.WithLockedNoResult(&n.requestsLock, func() {
		n.requests = utils.Filter(n.requests, func(r _Request, i int) bool {
			if r.key == key {
				go r.from.Receive(key, val)
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
		go from.Receive(key, val)
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

func NewBaseNode(net *netp.Network[BaseNode]) *BaseNode {
	return &BaseNode{store: make(map[string]string), network: net}
}
