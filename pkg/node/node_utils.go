package node

import (
	"fmt"
	"math/rand"
	"time"
)

func (n *Node) AddToStore(key, val string) {
	n.storageLock.Lock()
	defer n.storageLock.Unlock()

	n.storage[key] = val
}

func (n *Node) HasInStore(key string) (string, bool) {
	n.storageLock.RLock()
	defer n.storageLock.RUnlock()

	v, ok := n.storage[key]
	return v, ok
}

func (n *Node) ReceivedKey(key string) (string, bool) {
	n.receivedValuesLock.RLock()
	defer n.receivedValuesLock.RUnlock()

	v, ok := n.receivedValues[key]
	return v, ok
}

func (n *Node) PutKey(key, val string) {
	n.receivedValuesLock.Lock()
	defer n.receivedValuesLock.Unlock()

	n.receivedValues[key] = val
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

func (n *Node) findRequest(mId int) (Request, bool) {

	for _, d := range n.routeRequests {
		if d.id == mId {
			return d, true
		}
	}

	return Request{}, false
}

func (n *Node) isInRequests(id int) bool {
	_, ok := n.findRequest(id)
	return ok
}

func (n *Node) setAwaitingFromForRequest(id int) (Request, bool) {
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
func (n *Node) setRouteForRequest(id int, to INode) (Request, bool) {
	for i := 0; i < len(n.routeRequests); i++ {
		if id == n.routeRequests[i].id {
			n.routeRequests[i].routedTo = append(n.routeRequests[i].routedTo, to)
			return n.routeRequests[i], len(n.routeRequests[i].routedTo) == 1 && n.routeRequests[i].awaitingFrom == nil
		}
	}
	panic("Setting routedTo for message that is not in store")
}

func (n *Node) addRequest(r Request) {
	n.routeRequests = append(n.routeRequests, r)
}
func (n *Node) removeRequest(id int, from INode) Request {
	for i := 0; i < len(n.routeRequests); i++ {
		if n.routeRequests[i].id == id {
			defer func() {
				n.routeRequests[i] = n.routeRequests[len(n.routeRequests)-1]
				n.routeRequests = n.routeRequests[:(len(n.routeRequests) - 1)]
			}()
			return n.routeRequests[i]
		}
	}
	errorMessage := fmt.Sprintf("%p tried to remove request that is not in store; message came from %p", n.GetSelfAddress(), from)
	panic(errorMessage)
}

func (n *Node) removeRequests(ids ...int) []Request {
	ans := []Request{}

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

func (n *Node) findDisruptedDownloads(failedNode INode) []Request {
	ans := []Request{}

	for _, each := range n.routeRequests {
		if each.awaitingFrom == failedNode {
			ans = append(ans, each)
		}
	}

	return ans
}

func (n *Node) findDisruptedDownloadsById(failedNode INode, ids []int) []Request {
	ans := []Request{}

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
func (n *Node) clearFromRoutedTo(failedNode INode) []Request {
	ans := []Request{}

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
func (n *Node) clearFromRoutedToById(failedNode INode, ids []int) []Request {
	ans := []Request{}

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

func matchFromAndId(disruptedRoutes []Request) map[INode][]int {
	ans := make(map[INode][]int)

	for _, r := range disruptedRoutes {
		ans[r.from] = append(ans[r.from], r.id)
	}

	return ans
}

func (n *Node) seeIfGoingToFail(method Method) bool {
	if n.getFailChance == nil {
		return false
	}

	return rand.Float64() < n.getFailChance(method)
}
