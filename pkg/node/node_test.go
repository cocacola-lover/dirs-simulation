package node

import (
	"testing"
	"time"
)

var newRouteRequest = func(id int, key string, from INode, routedTo []INode, sentTo []INode, awaitingFrom INode) _Request {
	return _Request{
		id:           id,
		key:          key,
		from:         from,
		routedTo:     routedTo,
		sentTo:       sentTo,
		awaitingFrom: awaitingFrom,
	}
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
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends, node2.getNetworkFriends = func() []INode { return []INode{node2} }, func() []INode { return []INode{node1} }

		node2.PutVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2}, []INode{node2}, node2))

		node1.ConfirmDownloadMessage(0, "value", node2)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node3} }
		node3.getNetworkFriends = func() []INode { return []INode{node2, node4} }
		node4.getNetworkFriends = func() []INode { return []INode{node3} }

		node4.PutVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2}, []INode{node2}, node2))
		node2.addRequest(newRouteRequest(0, "key", node1, []INode{node3}, []INode{node3}, node3))
		node3.addRequest(newRouteRequest(0, "key", node2, []INode{node4}, []INode{node4}, node4))

		node3.ConfirmDownloadMessage(0, "value", node4)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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

		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends, node2.getNetworkFriends = func() []INode { return []INode{node2} }, func() []INode { return []INode{node1} }

		node2.PutVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2}, []INode{node2}, nil))

		node2.ReceiveDownloadMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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

		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node3} }
		node3.getNetworkFriends = func() []INode { return []INode{node2, node4} }
		node4.getNetworkFriends = func() []INode { return []INode{node3} }

		node4.PutVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2}, []INode{node2}, nil))
		node2.addRequest(newRouteRequest(0, "key", node1, []INode{node3}, []INode{node3}, nil))
		node3.addRequest(newRouteRequest(0, "key", node2, []INode{node4}, []INode{node4}, nil))

		node2.ReceiveDownloadMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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

func TestNode_ConfirmRouteMessage(t *testing.T) {

	t.Run("Test simple ConfirmRouteMessage", func(t *testing.T) {

		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends, node2.getNetworkFriends = func() []INode { return []INode{node2} }, func() []INode { return []INode{node1} }

		node2.PutVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, nil, []INode{node2}, nil))
		node1.ConfirmRouteMessage(0, node2)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node3} }
		node3.getNetworkFriends = func() []INode { return []INode{node2, node4} }
		node4.getNetworkFriends = func() []INode { return []INode{node3} }

		node4.PutVal("key", "value")

		node1.addRequest(newRouteRequest(0, "key", node1, nil, []INode{node2}, nil))
		node2.addRequest(newRouteRequest(0, "key", node1, nil, []INode{node3}, nil))
		node3.addRequest(newRouteRequest(0, "key", node2, nil, []INode{node4}, nil))

		node3.ConfirmRouteMessage(0, node4)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends, node2.getNetworkFriends = func() []INode { return []INode{node2} }, func() []INode { return []INode{node1} }

		node2.PutVal("key", "value")

		node1.ReceiveRouteMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

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
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node3} }
		node3.getNetworkFriends = func() []INode { return []INode{node2, node4} }
		node4.getNetworkFriends = func() []INode { return []INode{node3} }

		node4.PutVal("key", "value")

		node1.ReceiveRouteMessage(0, "key", node1)

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 20, 32)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})
}

func TestNode_Fail(t *testing.T) {

	t.Run("Test change routedTo on fail", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 3 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 3 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 3 })

		node1.getNetworkFriends = func() []INode { return []INode{node2, node3} }
		node2.getNetworkFriends = func() []INode { return []INode{node1} }
		node3.getNetworkFriends = func() []INode { return []INode{node1} }

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2, node3}, []INode{node2, node3}, nil))

		node3.Fail()

		ok, error := failAtButSuccedAt(func() bool {
			return len(node1.routeRequests[0].routedTo) == 1
		}, 1, 5)

		if !ok {
			t.Fatalf("Failed did not forced reroute : %s\n", error)
		}
	})

	t.Run("Test change download target on fail", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2, node3} }
		node2.getNetworkFriends = func() []INode { return []INode{node1} }
		node3.getNetworkFriends = func() []INode { return []INode{node1} }

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2, node3}, []INode{node2, node3}, node2))

		node2.PutVal("key", "value")
		node3.PutVal("key", "value")

		node2.Fail()

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 5, 10)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

	t.Run("Test change routedTo on fail through neighbour", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2, node3} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node4} }
		node3.getNetworkFriends = func() []INode { return []INode{node1} }
		node4.getNetworkFriends = func() []INode { return []INode{node2} }

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node2, node3}, []INode{node2, node3}, nil))
		node2.addRequest(newRouteRequest(0, "key", node1, []INode{node4}, []INode{node4}, nil))

		node4.Fail()

		ok, error := failAtButSuccedAt(func() bool {
			return len(node1.routeRequests[0].routedTo) == 1
		}, 1, 5)

		if !ok {
			t.Fatalf("Failed did not forced reroute : %s\n", error)
		}
	})

	t.Run("Test change download target on fail through neighbour", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2, node3} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node4} }
		node3.getNetworkFriends = func() []INode { return []INode{node1} }
		node4.getNetworkFriends = func() []INode { return []INode{node2} }

		// for i, n := range []INode{node1, node2, node3, node4} {
		// 	t.Logf("Node%d has a pointer of %p\n", i+1, n)
		// }

		node1.addRequest(newRouteRequest(0, "key", node1, []INode{node3}, []INode{node2, node3}, node2))
		node2.addRequest(newRouteRequest(0, "key", node1, []INode{}, []INode{node4}, node4))

		node3.PutVal("key", "value")
		node4.PutVal("key", "value")

		node4.Fail()

		ok, error := failAtButSuccedAt(func() bool {
			value, ok := node1.HasKey("key")

			if !ok || value != "value" {
				return false
			}
			return true
		}, 5, 11)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

}

func TestNode(t *testing.T) {

	t.Run("Test end-to-end two messages", func(t *testing.T) {
		node1 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node2 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node3 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })
		node4 := NewNode(1, 1, nil, func(with INode) (int, int) { return 1, 1 })

		node1.getNetworkFriends = func() []INode { return []INode{node2} }
		node2.getNetworkFriends = func() []INode { return []INode{node1, node3} }
		node3.getNetworkFriends = func() []INode { return []INode{node2, node4} }
		node4.getNetworkFriends = func() []INode { return []INode{node3} }

		node4.PutVal("key", "value")

		node2.StartSearch("key")
		time.Sleep(time.Millisecond)
		node1.StartSearch("key")

		ok, error := failAtButSuccedAt(func() bool {
			v1, ok1 := node1.HasKey("key")
			v2, ok2 := node2.HasKey("key")

			if !ok1 || !ok2 || v1 != "value" || v2 != "value" {
				return false
			}
			return true
		}, 20, 35)

		if !ok {
			t.Fatalf("Adding to store failed : %s\n", error)
		}
	})

}
