package basenode

import (
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

	friends []*BaseNode
}

// There is NO LOCK for Friends setter
func (n *BaseNode) SetFriends(friends []*BaseNode) {
	n.friends = friends
}

func (n *BaseNode) Receive(key string, val string) {

	utils.WithLockedNoResult(&n.storeLock, func() {
		n.store[key] = val
	})

	for i, r := range n.requests {
		if r.key == key {
			go r.from.Receive(key, val)
			utils.WithLockedNoResult(&n.requestsLock, func() {
				// Remove from array
				n.requests[i] = n.requests[len(n.requests)-1]
				n.requests = n.requests[:len(n.requests)-1]
			})
		}
	}
}

func (n *BaseNode) Ask(key string, from *BaseNode) {
	val, ok := n.store[key]

	if ok {
		go from.Receive(key, val)
	} else {
		utils.WithLockedNoResult(&n.requestsLock, func() {
			n.requests = append(n.requests, _Request{
				from: from, key: key,
			})
		})
	}
}
