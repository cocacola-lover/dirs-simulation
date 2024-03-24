package node

import (
	bmp "dirs/simulation/pkg/bandwidthManager"
	"sync"
	"time"
)

type Node struct {
	// Store holds key-value pairs
	store     map[string]string
	storeLock sync.RWMutex

	routeRequests     []_RouteRequest
	routeRequestsLock sync.RWMutex

	doneMessages     map[int]bool
	doneMessagesLock sync.RWMutex

	bm          *bmp.BandwidthManager
	selfAddress INode

	getNetworkFriends func() []INode
	getNetworkTunnel  func(with INode) (int, int)
	// OuterGetterFunctions^
}

func (n *Node) ReceiveRouteMessage(id int, key string, from INode) bool {

	if n.selfAddress != from {
		_, tunnelLength := n.getTunnel(from)
		time.Sleep(time.Millisecond * time.Duration(tunnelLength))
	}

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInRequests(id) || n.isInDoneMessages(id) {
		return false
	}

	_, ok := n.HasKey(key)

	if ok {
		go from.ConfirmRouteMessage(id, n.selfAddress)
	} else {
		r := _NewRouteRequest(id, key, from)
		r.sentTo = n.forEachFriendExcept(func(f INode) {
			go f.ReceiveRouteMessage(id, key, n.selfAddress)
		}, from)

		n.addRequest(r)
	}
	return true
}

func (n *Node) ConfirmRouteMessage(id int, from INode) bool {

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) || n.isInConfirmedRequests(id) {
		return false
	}

	r := n.setRouteForRequest(id, from)

	if r.from == n.selfAddress {
		go from.ReceiveDownloadMessage(id, r.key, n.selfAddress)
	} else {
		go r.from.ConfirmRouteMessage(id, n.selfAddress)
	}

	for _, nn := range r.sentTo {
		if nn != from {
			go nn.TimeoutRouteMessage(id, n.selfAddress)
		}
	}
	return true
}

func (n *Node) TimeoutRouteMessage(id int, from INode) bool {

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) {
		return false
	}

	n.doneMessagesLock.Lock()
	defer n.doneMessagesLock.Unlock()

	r, ok := n.findRequest(id)

	if !ok || r.from != from {
		return false
	} else {
		n.doneMessages[id] = true
		n.removeRequest(id)

		for _, nn := range r.sentTo {
			go nn.TimeoutRouteMessage(id, n.selfAddress)
		}
	}

	return true
}

func (n *Node) ReceiveDownloadMessage(id int, key string, from INode) {
	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	val, ok := n.HasKey(key)

	if ok {
		go from.ConfirmDownloadMessage(id, val, n.selfAddress)
	} else {
		n.routeRequestsLock.Lock()
		defer n.routeRequestsLock.Unlock()

		r, _ := n.findRequest(id)
		go r.routedTo.ReceiveDownloadMessage(id, key, n.selfAddress)
	}
}

func (n *Node) ConfirmDownloadMessage(id int, val string, from INode) {
	tunnelWidth, tunnelLength := n.getTunnel(from)
	n.bm.RegisterDownload(len(val), from.Bm(), tunnelWidth, tunnelLength, func(_ int) {
		n.doneMessagesLock.Lock()
		n.doneMessages[id] = true
		n.doneMessagesLock.Unlock()

		n.routeRequestsLock.Lock()
		defer n.routeRequestsLock.Unlock()

		r := n.removeRequest(id)

		if r.from == n.selfAddress {
			n.PutVal(r.key, val)
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
