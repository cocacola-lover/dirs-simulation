package basenode

func (n *BaseNode) GetFromStore(key string) (string, bool) {
	n.storeLock.RLock()
	defer n.storeLock.RUnlock()

	val, ok := n.store[key]
	return val, ok
}

func (n *BaseNode) PutInStore(m IMessage, val string) {
	if n.watchPutInStore != nil {
		go n.watchPutInStore(m, val)
	}

	n.storeLock.Lock()
	defer n.storeLock.Unlock()

	n.store[m.Key()] = val
}

func (n *BaseNode) AddRequest(ms ...IMessage) {
	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = append(n.requests, ms...)
}

func (n *BaseNode) HasMessage(m IMessage) bool {
	n.requestsLock.RLock()
	defer n.requestsLock.RUnlock()

	for _, mc := range n.requests {
		if mc.Id() == m.Id() {
			return true
		}
	}

	return false
}

func (n *BaseNode) RegisterDownload(m IMessage, val string) {

	if n.watchRegisterDownload != nil {
		n.watchRegisterDownload(m, val)
	}

	var tunnelWidth, tunnelLength int
	if n.getTunnel == nil {
		tunnelWidth = 1
		tunnelLength = 0
	} else {
		tunnelWidth, tunnelLength = n.getTunnel(m.From())
	}

	go n.bandwidthManager.RegisterDownload(len(val), m.From().bandwidthManager, tunnelWidth, tunnelLength, func() {
		m.From().Receive(m, val)
	})
}

// Assumes getFriends returns array with unique values +
// order of friends does not matter.
func (n *BaseNode) WhoToAsk(m IMessage) []*BaseNode {
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
