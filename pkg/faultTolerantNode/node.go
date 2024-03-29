package faulttolerantnode

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	"sync"
	"time"
)

type Node struct {
	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex

	routeRequests     []_Request
	routeRequestsLock sync.RWMutex

	doneMessages     map[int]bool
	doneMessagesLock sync.RWMutex

	bm          *bmp.BandwidthManager
	selfAddress *Node

	getNetworkFriends func() []*Node
	getNetworkTunnel  func(with *Node) (int, int)
	// OuterGetterFunctions^
}

func (n *Node) ReceiveRouteMessage(id int, key string, from *Node) bool {
	if n.selfAddress != from {
		_, tunnelLength := n.getTunnel(from)
		time.Sleep(time.Millisecond * time.Duration(tunnelLength))
	}

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
		r.sentTo = n.forEachFriendExcept(func(f *Node) {
			go f.ReceiveRouteMessage(id, key, n.selfAddress)
		}, from)

		n.addRequest(r)
	}
	return true
}

func (n *Node) ConfirmRouteMessage(id int, from *Node) bool {

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) {
		return false
	}

	r := n.setRouteForRequest(id, from)

	if r.from == n.selfAddress {
		go from.ReceiveDownloadMessage(id, r.key, n.selfAddress)
	} else {
		go r.from.ConfirmRouteMessage(id, n.selfAddress)
	}

	return true
}

// Nil "about" means that neighbor node failed, otherwise
// failure was experienced down the route.
func (n *Node) ReceiveFaultMessage(from *Node, about []int) {
	n.waitWayFrom(from)

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if about == nil {
		disruptedRoutes := []_Request{}

		for _, eachR := range n.findDisruptedDownloads(from) {
			if recoveredR, ok := n.setAwaitingFromForRequest(eachR.id); ok {
				go recoveredR.awaitingFrom.ReceiveDownloadMessage(eachR.id, eachR.key, n.selfAddress)
			} else {
				disruptedRoutes = append(disruptedRoutes, recoveredR)
			}
		}

		for _, emptyR := range n.clearFromRoutedTo(from) {
			
		}
	}

}

func (n *Node) ReceiveDownloadMessage(id int, key string, from *Node) {
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

func (n *Node) ConfirmDownloadMessage(id int, val string, from *Node) {
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
func (n *Node) SetSelfAddress(nn *Node) {
	n.selfAddress = nn
}
func (n *Node) GetSelfAddress() *Node {
	return n.selfAddress
}

// func NewNode(maxDownload int, maxUpload int, getNetworkFriends func() []INode, getNetworkTunnel func(with INode) (int, int)) *Node {
// 	n := &Node{
// 		bm:                bmp.NewBandwidthManager(maxDownload, maxUpload),
// 		store:             make(map[string]string),
// 		doneMessages:      make(map[int]bool),
// 		getNetworkFriends: getNetworkFriends,
// 		getNetworkTunnel:  getNetworkTunnel,
// 	}
// 	n.selfAddress = n
// 	return n
// }

// func (n *Node) SetOuterFunctions(getNetworkFriends func() []INode, getNetworkTunnel func(with INode) (int, int)) *Node {
// 	n.getNetworkFriends = getNetworkFriends
// 	n.getNetworkTunnel = getNetworkTunnel
// 	return n
