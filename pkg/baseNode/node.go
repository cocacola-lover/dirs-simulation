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

	removalCounter := 0
	for i, r := range n.requests {
		if r.key == key {
			go r.from.Receive(key, val)
			utils.WithLockedNoResult(&n.requestsLock, func() {
				// Remove from array
				n.requests[i-removalCounter] = n.requests[len(n.requests)-1]
				n.requests = n.requests[:len(n.requests)-1]

				removalCounter++
			})
		}
	}

}

func (n *BaseNode) Ask(key string, from *BaseNode) {
	val, ok := n.store[key]

	if ok {
		go from.Receive(key, val)
	} else {

		if len(n.friends) == 0 {
			return
		}

		if len(n.friends) == 1 && n.friends[0] == from {
			return
		}

		utils.WithLockedNoResult(&n.requestsLock, func() {
			n.requests = append(n.requests, _Request{
				from: from, key: key,
			})
		})

		for _, friend := range n.friends {
			if friend == from {
				continue
			}
			go friend.Ask(key, n)
		}

	}
}

func NewBaseNode() *BaseNode {
	return &BaseNode{store: make(map[string]string)}
}
