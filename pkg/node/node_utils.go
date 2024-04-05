package node

import (
	"time"
)

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

func (n *Node) findRequest(mId int) (_Request, bool) {

	for _, d := range n.routeRequests {
		if d.id == mId {
			return d, true
		}
	}

	return _Request{}, false
}

func (n *Node) isInRequests(id int) bool {
	_, ok := n.findRequest(id)
	return ok
}

func (n *Node) setAwaitingFromForRequest(id int) (_Request, bool) {
	for i := 0; i < len(n.routeRequests); i++ {
		if id == n.routeRequests[i].id {
			if len(n.routeRequests[i].routedTo) == 0 {
				return n.routeRequests[i], false
			}
			// Remove first from routedTo
			defer func() {
				n.routeRequests[i].routedTo[0] = n.routeRequests[i].routedTo[len(n.routeRequests[i].routedTo)-1]
				n.routeRequests[i].routedTo = n.routeRequests[i].routedTo[:len(n.routeRequests[i].routedTo)-1]
			}()

			n.routeRequests[i].awaitingFrom = n.routeRequests[i].routedTo[0]
			return n.routeRequests[i], true
		}
	}
	panic("Setting awaitingFor for message that is not in store")
}

// The second return value tells if it's first routedTo value
func (n *Node) setRouteForRequest(id int, to INode) (_Request, bool) {
	for i := 0; i < len(n.routeRequests); i++ {
		if id == n.routeRequests[i].id {
			n.routeRequests[i].routedTo = append(n.routeRequests[i].routedTo, to)
			return n.routeRequests[i], len(n.routeRequests[i].routedTo) == 1 && n.routeRequests[i].awaitingFrom == nil
		}
	}
	panic("Setting routedTo for message that is not in store")
}

func (n *Node) addRequest(r _Request) {
	n.routeRequests = append(n.routeRequests, r)
}
func (n *Node) removeRequest(id int) _Request {
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

func (n *Node) removeRequests(ids ...int) []_Request {
	ans := []_Request{}

	for i := 0; i < len(n.routeRequests); i++ {
		for _, id := range ids {
			if n.routeRequests[i].id == id {
				ans = append(ans, n.routeRequests[i])
				n.routeRequests[i] = n.routeRequests[len(n.routeRequests)-1]
				n.routeRequests = n.routeRequests[:(len(n.routeRequests) - 1)]
			}
		}
	}

	if len(ans) != len(ids) {
		panic("Tried to remove requests that is not in store")
	}

	return ans
}

func (n *Node) isInDoneMessages(id int) bool {
	n.doneMessagesLock.RLock()
	defer n.doneMessagesLock.RUnlock()

	_, ok := n.doneMessages[id]
	return ok
}

func (n *Node) waitWayFrom(from INode) {
	if n.selfAddress != from {
		_, tunnelLength := n.getTunnel(from)
		time.Sleep(time.Millisecond * time.Duration(tunnelLength))
	}
}

func (n *Node) findDisruptedDownloads(failedNode INode) []_Request {
	ans := []_Request{}

	for _, each := range n.routeRequests {
		if each.awaitingFrom == failedNode {
			ans = append(ans, each)
		}
	}

	return ans
}

func (n *Node) findDisruptedDownloadsById(failedNode INode, ids []int) []_Request {
	ans := []_Request{}

	for _, each := range n.routeRequests {
		if each.awaitingFrom == failedNode {
			for _, id := range ids {
				if each.id == id {
					ans = append(ans, each)
					break
				}
			}
		}
	}

	return ans
}

// Returns requests that are left with empty routedTo
func (n *Node) clearFromRoutedTo(failedNode INode) []_Request {
	ans := []_Request{}

	for i := range n.routeRequests {
		for j, eachRoutedTo := range n.routeRequests[i].routedTo {
			if eachRoutedTo == failedNode {
				n.routeRequests[i].routedTo[j] = n.routeRequests[i].routedTo[len(n.routeRequests[i].routedTo)-1]
				n.routeRequests[i].routedTo = n.routeRequests[i].routedTo[:len(n.routeRequests[i].routedTo)-1]

				if len(n.routeRequests[i].routedTo) == 0 {
					ans = append(ans, n.routeRequests[i])
					break
				}
			}
		}
	}

	return ans
}

// Returns requests that are left with empty routedTo
func (n *Node) clearFromRoutedToById(failedNode INode, ids []int) []_Request {
	ans := []_Request{}

	for i := range n.routeRequests {
		for _, id := range ids {
			if n.routeRequests[i].id == id {
				for j, eachRoutedTo := range n.routeRequests[i].routedTo {
					if eachRoutedTo == failedNode {
						n.routeRequests[i].routedTo[j] = n.routeRequests[i].routedTo[len(n.routeRequests[i].routedTo)-1]
						n.routeRequests[i].routedTo = n.routeRequests[i].routedTo[:len(n.routeRequests[i].routedTo)-1]

						if len(n.routeRequests[i].routedTo) == 0 {
							ans = append(ans, n.routeRequests[i])
							break
						}
					}
				}
			}
		}
	}

	return ans
}

func matchFromAndId(disruptedRoutes []_Request) map[INode][]int {
	ans := make(map[INode][]int)

	for _, r := range disruptedRoutes {
		ans[r.from] = append(ans[r.from], r.id)
	}

	return ans
}
