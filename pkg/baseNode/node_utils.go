package basenode

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

func (n *BaseNode) AddRequest(key string, from *BaseNode) {
	n.requestsLock.Lock()
	defer n.requestsLock.Unlock()

	n.requests = append(n.requests, _Request{
		from: from, key: key,
	})
}

func (n *BaseNode) RegisterDownload(key string, val string, with *BaseNode) {
	n.bandwidthManager.RegisterDownload(
		len(val),
		with.bandwidthManager,
		n.network.GetTunnel(n, with),
		func() {
			with.Receive(key, val)
		},
	)
}
