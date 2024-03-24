package node

func (n *Node) HasKey(key string) (string, bool) {
	n.storeLock.RLock()
	defer n.storeLock.RUnlock()

	v, ok := n.store[key]
	return v, ok
}

func (n *Node) PutVal(key, val string) {
	n.storeLock.Lock()
	defer n.storeLock.Unlock()

	n.store[key] = val
}

// Returns an array of friends to which did "do"
func (n *Node) forEachFriendExcept(do func(f INode), except INode) []INode {
	if n.getNetworkFriends == nil {
		return nil
	}

	fs := n.getNetworkFriends()

	for i := 0; i < len(fs); i++ {
		if fs[i] == except {
			fs[i] = fs[len(fs)-1]
			fs = fs[:len(fs)-1]

			if i == len(fs) {
				break
			}
		}
		do(fs[i])
	}

	return fs
}

// TunnelWidth and TunnelLength
func (n *Node) getTunnel(to INode) (int, int) {
	var tunnelWidth, tunnelLength int
	if n.getNetworkTunnel == nil {
		tunnelWidth = 1
		tunnelLength = 0
	} else {
		tunnelWidth, tunnelLength = n.getNetworkTunnel(to)
	}

	return tunnelWidth, tunnelLength
}

func (n *Node) findRequest(mId int) (_RouteRequest, bool) {

	for _, d := range n.routeRequests {
		if d.id == mId {
			return d, true
		}
	}

	return _RouteRequest{}, false
}

func (n *Node) isInRequests(id int) bool {
	_, ok := n.findRequest(id)
	return ok
}
func (n *Node) isInConfirmedRequests(id int) bool {
	r, ok := n.findRequest(id)
	if !ok {
		return false
	}
	return r.routedTo != nil
}

func (n *Node) setRouteForRequest(id int, to INode) _RouteRequest {
	for i := 0; i < len(n.routeRequests); i++ {
		if id == n.routeRequests[i].id {
			n.routeRequests[i].routedTo = to
			return n.routeRequests[i]
		}
	}
	panic("Setting routedTo for message that is not in store")
}

func (n *Node) addRequest(r _RouteRequest) {
	n.routeRequests = append(n.routeRequests, r)
}
func (n *Node) removeRequest(id int) _RouteRequest {
	for i := 0; i < len(n.routeRequests); i++ {
		if n.routeRequests[i].id == id {
			defer func() {
				n.routeRequests[i] = n.routeRequests[len(n.routeRequests)-1]
				n.routeRequests = n.routeRequests[:(len(n.routeRequests) - 1)]
			}()
			return n.routeRequests[i]
		}
	}
	panic("Tried to remove request that is not in store")
}

func (n *Node) isInDoneMessages(id int) bool {
	n.doneMessagesLock.RLock()
	defer n.doneMessagesLock.RUnlock()

	_, ok := n.doneMessages[id]
	return ok
}
