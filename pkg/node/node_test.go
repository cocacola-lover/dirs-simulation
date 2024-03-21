package node

import (
	"testing"
	"time"
)

var friendsFactory = func(nodes ...INode) func() []INode {
	return func() []INode {
		return nodes
	}
}

var tunnelFactory = func() func(with INode) (int, int) {
	return func(with INode) (int, int) {
		return 1, 1
	}
}

var newRouteRequest = func(id int, key string, from INode, routedTo INode, sentTo []INode) _RouteRequest {
	return _RouteRequest{id: id, key: key, from: from, routedTo: routedTo, sentTo: sentTo}
}

var failAtButSuccedAt = func(testFu func() bool, before int, after int) (bool, string) {
	time.Sleep(time.Millisecond * time.Duration(before))
	if testFu() {
		return false, "Succeded too early"
	}

	time.Sleep(time.Millisecond * time.Duration(after-before))
	if !testFu() {
		return false, "Failed too late"
	}

	return true, ""
}

func TestNode_ConfirmDownloadMessage(t *testing.T) {
	t.Run("Test single ConfirmDownloadMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node1.getNetworkFriends, node2.getNetworkFriends = friendsFactory(node2), friendsFactory(node1)

		node2.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, node2, nil))

		node1.ConfirmDownloadMessage(0, "value", node2)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 4, 8)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

	t.Run("Test chain ConfirmDownloadMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node3 := NewNode(1, 1, nil, tunnelFactory())
		node4 := NewNode(1, 1, nil, tunnelFactory())

		node1.getNetworkFriends = friendsFactory(node2)
		node2.getNetworkFriends = friendsFactory(node1, node3)
		node3.getNetworkFriends = friendsFactory(node2, node4)
		node4.getNetworkFriends = friendsFactory(node3)

		node4.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, node2, nil))
		node2.addRequest(newRouteRequest(0, "key", node1, node3, nil))
		node3.addRequest(newRouteRequest(0, "key", node2, node4, nil))

		node3.ConfirmDownloadMessage(0, "value", node4)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 15, 20)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})
}

func TestNode_ReceiveDownloadMessage(t *testing.T) {

	t.Run("Test simple ReceiveDownloadMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node1.getNetworkFriends, node2.getNetworkFriends = friendsFactory(node2), friendsFactory(node1)

		node2.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, node2, nil))

		node2.ReceiveDownloadMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 3, 7)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

	t.Run("Test chain ReceiveDownloadMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node3 := NewNode(1, 1, nil, tunnelFactory())
		node4 := NewNode(1, 1, nil, tunnelFactory())

		node1.getNetworkFriends = friendsFactory(node2)
		node2.getNetworkFriends = friendsFactory(node1, node3)
		node3.getNetworkFriends = friendsFactory(node2, node4)
		node4.getNetworkFriends = friendsFactory(node3)

		node4.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, node2, nil))
		node2.addRequest(newRouteRequest(0, "key", node1, node3, nil))
		node3.addRequest(newRouteRequest(0, "key", node2, node4, nil))

		node2.ReceiveDownloadMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 15, 24)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})
}

func TestNode_TimeoutRouteMessage(t *testing.T) {

	t.Run("Test simple TimeoutRouteMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node3 := NewNode(1, 1, nil, tunnelFactory())
		node1.getNetworkFriends, node2.getNetworkFriends = friendsFactory(node2), friendsFactory(node1)

		node2.addRequest(newRouteRequest(0, "key", node1, node3, nil))
		node2.TimeoutRouteMessage(0, node1)

		time.Sleep(2 * time.Millisecond)

		node2.routeRequestsLock.Lock()
		_, ok := node2.findRequest(0)
		node2.routeRequestsLock.Unlock()

		if ok {
			t.Fatal("Timeout did not delete messages")
		}
	})

	t.Run("Test chain TimeoutRouteMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node3 := NewNode(1, 1, nil, tunnelFactory())
		node4 := NewNode(1, 1, nil, tunnelFactory())
		node5 := NewNode(1, 1, nil, tunnelFactory())

		node1.getNetworkFriends = friendsFactory(node2)
		node2.getNetworkFriends = friendsFactory(node1, node3)
		node3.getNetworkFriends = friendsFactory(node2, node4)
		node4.getNetworkFriends = friendsFactory(node3, node5)

		node4.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, node2, []INode{node2}))
		node2.addRequest(newRouteRequest(0, "key", node1, node3, []INode{node3}))
		node3.addRequest(newRouteRequest(0, "key", node2, node4, []INode{node4}))
		node4.addRequest(newRouteRequest(0, "key", node3, node5, []INode{node5}))

		node2.TimeoutRouteMessage(0, node1)

		ok, error := failAtButSuccedAt(func() bool {
			node2.routeRequestsLock.Lock()
			_, ok2 := node2.findRequest(0)
			node2.routeRequestsLock.Unlock()

			node3.routeRequestsLock.Lock()
			_, ok3 := node3.findRequest(0)
			node3.routeRequestsLock.Unlock()

			node4.routeRequestsLock.Lock()
			_, ok4 := node4.findRequest(0)
			node4.routeRequestsLock.Unlock()

			return !ok2 && !ok3 && !ok4
		}, 1, 6)

		if !ok {
			t.Fatalf("Timeout did not delete messages : %s\n", error)
		}
	})
}

func TestNode_ConfirmRouteMessage(t *testing.T) {

	t.Run("Test simple ConfirmRouteMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node1.getNetworkFriends, node2.getNetworkFriends = friendsFactory(node2), friendsFactory(node1)
		node2.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, nil, nil))
		node1.ConfirmRouteMessage(0, node2)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 5, 8)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

	t.Run("Test chain ConfirmRouteMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node3 := NewNode(1, 1, nil, tunnelFactory())
		node4 := NewNode(1, 1, nil, tunnelFactory())

		node1.getNetworkFriends = friendsFactory(node2)
		node2.getNetworkFriends = friendsFactory(node1, node3)
		node3.getNetworkFriends = friendsFactory(node2, node4)
		node4.getNetworkFriends = friendsFactory(node3)

		node4.putVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, nil, []INode{node2}))
		node2.addRequest(newRouteRequest(0, "key", node1, nil, []INode{node3}))
		node3.addRequest(newRouteRequest(0, "key", node2, nil, []INode{node4}))

		node3.ConfirmRouteMessage(0, node4)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 20, 26)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})
}

func TestNode_ReceiveRouteMessage(t *testing.T) {

	t.Run("Test simple ReceiveRouteMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node1.getNetworkFriends, node2.getNetworkFriends = friendsFactory(node2), friendsFactory(node1)
		node2.putVal("key", "value")

		// node1.addRequest(newRouteRequest(0, "key", node1, nil, nil))
		node1.ReceiveRouteMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 5, 11)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

	t.Run("Test chain ReceiveRouteMessage", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, tunnelFactory())
		node2 := NewNode(1, 1, nil, tunnelFactory())
		node3 := NewNode(1, 1, nil, tunnelFactory())
		node4 := NewNode(1, 1, nil, tunnelFactory())

		node1.getNetworkFriends = friendsFactory(node2)
		node2.getNetworkFriends = friendsFactory(node1, node3)
		node3.getNetworkFriends = friendsFactory(node2, node4)
		node4.getNetworkFriends = friendsFactory(node3)

		node4.putVal("key", "value")

		node1.ReceiveRouteMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.hasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 20, 30)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})
}
