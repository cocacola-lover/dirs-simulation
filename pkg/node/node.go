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

	bm *bmp.BandwidthManager

	getNetworkFriends func() []INode
	getNetworkTunnel  func(with INode) (int, int)
	// OuterGetterFunctions^
}

func (n *Node) ReceiveRouteMessage(id int, key string, from INode) {

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInRequests(id) || n.isInDoneMessages(id) {
		return
	}

	_, ok := n.hasKey(key)

	if ok {
		go from.ConfirmRouteMessage(id, n)
	} else {
		r := _NewRouteRequest(id, key, from)
		r.sentTo = n.forEachFriendExcept(func(f INode) {
			go f.ReceiveRouteMessage(id, key, n)
		}, from)

		n.addRequest(r)
	}
}

func (n *Node) ConfirmRouteMessage(id int, from INode) {

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) || n.isInConfirmedRequests(id) {
		return
	}

	r := n.setRouteForRequest(id, from)

	if n.requestCameFromMe(id) {
		go from.ReceiveDownloadMessage(id, r.key, n)
	} else {
		go r.from.ConfirmRouteMessage(id, n)
	}

	for _, nn := range r.sentTo {
		if nn != from {
			go nn.TimeoutRouteMessage(id, n)
		}
	}
}

func (n *Node) TimeoutRouteMessage(id int, from INode) {

	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	n.routeRequestsLock.Lock()
	defer n.routeRequestsLock.Unlock()

	if n.isInDoneMessages(id) {
		return
	}

	n.doneMessagesLock.Lock()
	defer n.doneMessagesLock.Unlock()

	r, ok := n.findRequest(id)

	if !ok || r.from != from {
		return
	} else {
		n.doneMessages[id] = true
		n.removeRequest(id)

		for _, nn := range r.sentTo {
			go nn.TimeoutRouteMessage(id, n)
		}
	}
}

func (n *Node) ReceiveDownloadMessage(id int, key string, from INode) {
	_, tunnelLength := n.getTunnel(from)
	time.Sleep(time.Millisecond * time.Duration(tunnelLength))

	val, ok := n.hasKey(key)

	if ok {
		go from.ConfirmDownloadMessage(id, val, n)
	} else {
		n.routeRequestsLock.Lock()
		defer n.routeRequestsLock.Unlock()

		r, _ := n.findRequest(id)
		go r.routedTo.ReceiveDownloadMessage(id, key, n)
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

		if r.from == n {
			n.putVal(r.key, val)
		} else {
			go r.from.ConfirmDownloadMessage(id, val, n)
		}
	})
}

func (n *Node) Bm() *bmp.BandwidthManager {
	return n.bm
}

func NewNode(maxDownload int, maxUpload int, getNetworkFriends func() []INode, getNetworkTunnel func(with INode) (int, int)) *Node {
	return &Node{
		bm:                bmp.NewBandwidthManager(maxDownload, maxUpload),
		store:             make(map[string]string),
		doneMessages:      make(map[int]bool),
		getNetworkFriends: getNetworkFriends,
		getNetworkTunnel:  getNetworkTunnel,
	}
}
