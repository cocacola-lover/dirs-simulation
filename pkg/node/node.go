package node

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	idgenerator "dirs/simulation/pkg/idGenerator.go"
	"fmt"
	"sync"
	"sync/atomic"
)

type Node struct {
	hasFailed atomic.Bool

	// Store holds key-value pairs
	storage     map[string]string
	storageLock sync.RWMutex

	receivedValues     map[string]string
	receivedValuesLock sync.RWMutex

	routeRequests     []Request
	routeRequestsLock sync.RWMutex

	doneMessages     map[int]bool
	doneMessagesLock sync.RWMutex

	bm          *bmp.BandwidthManager
	selfAddress INode

	getNetworkFriends func() []INode
	getNetworkTunnel  func(with INode) (int, int)
	getFailChance     func(method Method) float64
	// OuterGetterFunctions^
}

func (n *Node) StartSearch(key string) int {
	if n.hasFailed.Load() {
		panic("Started search on failed node")
	}
	id := idgenerator.GetId()

	n.selfAddress.ReceiveRouteMessage(id, key, n.selfAddress)

	return id
}

func (n *Node) ReceiveRouteMessage(id int, key string, from INode) bool {
	if n.hasFailed.Load() {
		return false
	}
	if n.seeIfGoingToFail(ReceiveRouteMethod) {
		n.GetSelfAddress().Fail()
		return false
	}

	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInRequests(id) || n.isInDoneMessages(id) {
		return false
	}

	_, ok := n.selfAddress.HasInStore(key)

	if ok {
		go from.ConfirmRouteMessage(id, n.selfAddress)
	} else {
		r := _NewRequest(id, key, from)
		r.sentTo = n.forEachFriendExcept(func(f INode) {
			go f.ReceiveRouteMessage(id, key, n.selfAddress)
		}, from)

		n.addRequest(r)
	}
	return true
}

func (n *Node) ConfirmRouteMessage(id int, from INode) bool {
	if n.hasFailed.Load() {
		return false
	}
	if n.seeIfGoingToFail(ConfirmRouteMethod) {
		n.GetSelfAddress().Fail()
		return false
	}
	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) {
		return false
	}

	r, first := n.setRouteForRequest(id, from)

	if first {
		if r.from == n.selfAddress {
			r, ok := n.setAwaitingFromForRequest(r.id)
			if !ok {
				panic("Can not set awating on ConfirmRouteMessage")
			}
			go from.ReceiveDownloadMessage(id, r.key, n.selfAddress)
		} else {
			go r.from.ConfirmRouteMessage(id, n.selfAddress)
		}
	}

	return true
}

func (n *Node) Fail() {
	if n.hasFailed.Load() {
		return
	}
	n.hasFailed.Store(true)
	for _, eachFriend := range n.getNetworkFriends() {
		go eachFriend.ReceiveFaultMessage(n.selfAddress, nil)
	}
}

// Nil "about" means that neighbor node failed, otherwise
// failure was experienced down the route.
// Returns disrupted ids
func (n *Node) ReceiveFaultMessage(from INode, about []int) []int {
	if n.hasFailed.Load() {
		return nil
	}
	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	disruptedRoutes := []Request{}

	if about == nil {
		for _, eachR := range n.findDisruptedDownloads(from) {
			if recoveredR, ok := n.setAwaitingFromForRequest(eachR.id); ok {
				go recoveredR.awaitingFrom.ReceiveDownloadMessage(eachR.id, eachR.key, n.selfAddress)
			} else {
				disruptedRoutes = append(disruptedRoutes, recoveredR)
			}
		}

		disruptedRoutes = append(disruptedRoutes, n.clearFromRoutedTo(from)...)
	} else {
		for _, eachR := range n.findDisruptedDownloadsById(from, about) {
			if recoveredR, ok := n.setAwaitingFromForRequest(eachR.id); ok {
				go recoveredR.awaitingFrom.ReceiveDownloadMessage(eachR.id, eachR.key, n.selfAddress)
			} else {
				disruptedRoutes = append(disruptedRoutes, recoveredR)
			}
		}

		disruptedRoutes = append(disruptedRoutes, n.clearFromRoutedToById(from, about)...)
	}

	disruptedIds := make([]int, 0, len(disruptedRoutes))

	n.doneMessagesLock.Lock()
	defer n.doneMessagesLock.Unlock()

	for node, ids := range matchFromAndId(disruptedRoutes) {
		for _, id := range ids {
			n.doneMessages[id] = true
			disruptedIds = append(disruptedIds, id)
		}

		if node != n.selfAddress {
			go node.ReceiveFaultMessage(n.selfAddress, ids)
		} else {
			rs := n.removeRequests(ids...)
			go n.selfAddress.RetryMessages(rs)
		}
	}

	return disruptedIds

}

// Returns new ids for retried messages
func (n *Node) RetryMessages(rs []Request) []int {
	if n.hasFailed.Load() {
		return nil
	}

	newIds := make([]int, 0, len(rs))
	for range rs {
		newIds = append(newIds, idgenerator.GetId())
	}

	for i, eachR := range rs {
		go n.selfAddress.ReceiveRouteMessage(newIds[i], eachR.key, n.selfAddress)
	}

	return newIds
}

func (n *Node) ReceiveDownloadMessage(id int, key string, from INode) {
	if n.hasFailed.Load() {
		return
	}
	if n.seeIfGoingToFail(ReceiveDownloadMethod) {
		n.GetSelfAddress().Fail()
		return
	}

	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) {
		return
	}

	val, ok := n.selfAddress.HasInStore(key)

	if ok {
		go from.ConfirmDownloadMessage(id, val, n.selfAddress)
	} else {

		r, ok := n.setAwaitingFromForRequest(id)
		if !ok {
			errorMessage := fmt.Sprintf("Can not set awating on receiveDownloadMessage with node %p", n.GetSelfAddress())
			panic(errorMessage)
		}
		go r.awaitingFrom.ReceiveDownloadMessage(id, key, n.selfAddress)
	}
}

func (n *Node) ConfirmDownloadMessage(id int, val string, from INode) {
	if n.hasFailed.Load() {
		return
	}
	if n.seeIfGoingToFail(ConfirmDownloadMethod) {
		n.GetSelfAddress().Fail()
		return
	}

	tunnelWidth, tunnelLength := n.getTunnel(from)
	n.bm.RegisterDownload(len(val), from.Bm(), tunnelWidth, tunnelLength, func(_ int) {
		n.doneMessagesLock.Lock()
		if _, ok := n.doneMessages[id]; ok {
			return
		}
		n.doneMessages[id] = true
		n.doneMessagesLock.Unlock()

		n.routeRequestsLock.Lock()
		defer n.routeRequestsLock.Unlock()

		r := n.removeRequest(id, from)

		if r.from == n.selfAddress {
			n.selfAddress.PutKey(r.key, val)
		} else {
			go r.from.ConfirmDownloadMessage(id, val, n.selfAddress)
		}
	})
}

func (n *Node) Bm() *bmp.BandwidthManager {
	return n.bm
}
func (n *Node) SetSelfAddress(nn INode) {
	n.selfAddress = nn
}
func (n *Node) GetSelfAddress() INode {
	return n.selfAddress
}

func (n *Node) HasFailed() bool {
	return n.hasFailed.Load()
}

func (n *Node) Close() {
	n.hasFailed.Store(true)
	n.Bm().Close()
}

func NewNode(maxDownload int, maxUpload int, getNetworkFriends func() []INode, getNetworkTunnel func(with INode) (int, int)) *Node {
	n := &Node{
		bm:                bmp.NewBandwidthManager(maxDownload, maxUpload),
		storage:           make(map[string]string),
		receivedValues:    make(map[string]string),
		doneMessages:      make(map[int]bool),
		getNetworkFriends: getNetworkFriends,
		getNetworkTunnel:  getNetworkTunnel,
	}
	n.selfAddress = n
	return n
}

func (n *Node) SetOuterFunctions(
	getNetworkFriends func() []INode,
	getNetworkTunnel func(with INode) (int, int),
	getFailChance func(method Method) float64,
) *Node {
	n.getNetworkFriends = getNetworkFriends
	n.getNetworkTunnel = getNetworkTunnel
	n.getFailChance = getFailChance
	return n
}
