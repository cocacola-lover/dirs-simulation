package node

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	idgenerator "dirs/simulation/pkg/idGenerator.go"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	hasFailed atomic.Bool

	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex

	routeRequests     []_Request
	routeRequestsLock sync.RWMutex

	doneMessages     map[int]bool
	doneMessagesLock sync.RWMutex

	bm          *bmp.BandwidthManager
	selfAddress INode

	getNetworkFriends func() []INode
	getNetworkTunnel  func(with INode) (int, int)
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
	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInRequests(id) || n.isInDoneMessages(id) {
		return false
	}

	_, ok := n.selfAddress.HasKey(key)

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
	n.hasFailed.Store(false)
	for _, eachFriend := range n.getNetworkFriends() {
		go eachFriend.ReceiveFaultMessage(n.selfAddress, nil)
	}
}

// Nil "about" means that neighbor node failed, otherwise
// failure was experienced down the route.
func (n *Node) ReceiveFaultMessage(from INode, about []int) {
	if n.hasFailed.Load() {
		return
	}
	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	disruptedRoutes := []_Request{}

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

	for node, ids := range matchFromAndId(disruptedRoutes) {
		if node != n.selfAddress {
			go node.ReceiveFaultMessage(n.selfAddress, ids)
		} else {
			go n.selfAddress.RetryMessages(ids)
		}
	}

}

// Returns new ids for retried messages
func (n *Node) RetryMessages(ids []int) []int {
	if n.hasFailed.Load() {
		return nil
	}
	newIds := make([]int, len(ids))
	for i := range ids {
		newIds[i] = idgenerator.GetId()
	}

	n.doneMessagesLock.Lock()
	for _, id := range ids {
		n.doneMessages[id] = true
	}
	n.doneMessagesLock.Unlock()

	n.routeRequestsLock.RLock()
	defer n.routeRequestsLock.RUnlock()

	rs := n.removeRequests(ids...)
	for i, eachR := range rs {
		go n.selfAddress.ReceiveRouteMessage(newIds[i], eachR.key, n.selfAddress)
	}

	return newIds
}

func (n *Node) ReceiveDownloadMessage(id int, key string, from INode) {
	if n.hasFailed.Load() {
		return
	}

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	val, ok := n.selfAddress.HasKey(key)

	if ok {
		go from.ConfirmDownloadMessage(id, val, n.selfAddress)
	} else {
		n.routeRequestsLock.Lock()
		defer n.routeRequestsLock.Unlock()

		r, ok := n.setAwaitingFromForRequest(id)
		if !ok {
			panic("Can not set awating on receiveDownloadMessage")
		}
		go r.awaitingFrom.ReceiveDownloadMessage(id, key, n.selfAddress)
	}
}

func (n *Node) ConfirmDownloadMessage(id int, val string, from INode) {
	if n.hasFailed.Load() {
		return
	}
	tunnelWidth, tunnelLength := n.getTunnel(from)
	n.bm.RegisterDownload(len(val), from.Bm(), tunnelWidth, tunnelLength, func(_ int) {
		n.doneMessagesLock.Lock()
		n.doneMessages[id] = true
		n.doneMessagesLock.Unlock()

		n.routeRequestsLock.Lock()
		defer n.routeRequestsLock.Unlock()

		r := n.removeRequest(id)

		if r.from == n.selfAddress {
			n.selfAddress.PutVal(r.key, val)
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

func (n *Node) Close() {
	n.Bm().Close()
}

func NewNode(maxDownload int, maxUpload int, getNetworkFriends func() []INode, getNetworkTunnel func(with INode) (int, int)) *Node {
	n := &Node{
		bm:                bmp.NewBandwidthManager(maxDownload, maxUpload),
		store:             make(map[string]string),
		doneMessages:      make(map[int]bool),
		getNetworkFriends: getNetworkFriends,
		getNetworkTunnel:  getNetworkTunnel,
	}
	n.selfAddress = n
	return n
}

func (n *Node) SetOuterFunctions(getNetworkFriends func() []INode, getNetworkTunnel func(with INode) (int, int)) *Node {
	n.getNetworkFriends = getNetworkFriends
	n.getNetworkTunnel = getNetworkTunnel
	return n
}
