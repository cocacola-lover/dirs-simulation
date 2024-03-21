package basenode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	fp "dirs/simulation/pkg/fundamentals"
)

func (n *BaseNode) getFromStore(key string) (string, bool) {
	n.storeLock.RLock()
	defer n.storeLock.RUnlock()

	val, ok := n.store[key]
	return val, ok
}

func (n *BaseNode) addToStore(key, val string) {
	n.storeLock.Lock()
	defer n.storeLock.Unlock()

	n.store[key] = val
}

func (n *BaseNode) registerInStore(m fp.IMessage, val string) {
	if n.watchPutInStore != nil {
		go n.watchPutInStore(m, n)
	}

	n.addToStore(m.Key(), val)
}

func (n *BaseNode) addRequest(ms ..._Request) {
	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = append(n.requests, ms...)
}

func (n *BaseNode) hasMessage(m fp.IMessage) bool {
	n.requestsLock.RLock()
	defer n.requestsLock.RUnlock()

	for _, mc := range n.requests {
		if mc.Id() == m.Id() {
			return true
		}
	}

	return false
}

func (n *BaseNode) registerDownload(m fp.IMessage, val string) {

	if n.watchRegisterDownload != nil {
		n.watchRegisterDownload(m, n)
	}

	var tunnelWidth, tunnelLength int
	if n.getTunnel == nil {
		tunnelWidth = 1
		tunnelLength = 0
	} else {
		tunnelWidth, tunnelLength = n.getTunnel(m.From())
	}

	go n.bandwidthManager.RegisterDownload(len(val), m.From().BandwidthManager(), tunnelWidth, tunnelLength, func(id int) {
		m.From().Receive(m, val)
	})
}

// Assumes getFriends returns array with unique values +
// order of friends does not matter.
func (n *BaseNode) whoToAsk(m fp.IMessage) []fp.INode {
	if n.getFriends == nil {
		return nil
	} else {
		friends := n.getFriends()
		if len(friends) == 0 || (len(friends) == 1 && friends[0] == m.From()) {
			return nil
		}

		for i := 0; i < len(friends); i++ {
			if friends[i] == m.From() {
				friends[i] = friends[len(friends)-1]
				friends = friends[:len(friends)-1]
				break
			}
		}

		return friends

	}
}

func (n *BaseNode) BandwidthManager() *bmp.BandwidthManager {
	return n.bandwidthManager
}
