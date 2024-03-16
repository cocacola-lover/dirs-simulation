package basenode

import mp "dirs/simulation/pkg/message"

func (n *BaseNode) GetFromStore(key string) (string, bool) {
	n.storeLock.RLock()
	defer n.storeLock.RUnlock()

	val, ok := n.store[key]
	return val, ok
}

func (n *BaseNode) PutInStore(key string, val string) {
	n.storeLock.Lock()
	defer n.storeLock.Unlock()

	n.store[key] = val
}

func (n *BaseNode) AddRequest(messages ...mp.BaseMessage[BaseNode]) {
	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = append(n.requests, messages...)
}

func (n *BaseNode) HasMessage(message mp.BaseMessage[BaseNode]) bool {
	n.requestsLock.RLock()
	defer n.requestsLock.RUnlock()

	for _, m := range n.requests {
		if m.Id == message.Id {
			return true
		}
	}

	return false
}

func (n *BaseNode) RegisterDownload(message mp.BaseMessage[BaseNode], val string) {
	go n.bandwidthManager.RegisterDownload(
		len(val),
		message.From.bandwidthManager,
		n.network.GetTunnel(n, message.From),
		func() {
			message.From.Receive(message, val)
		},
	)
}
