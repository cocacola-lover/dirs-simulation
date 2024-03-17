package basenode

import (
	fp "dirs/simulation/pkg/fundamentals"
	"sync/atomic"
	"testing"
	"time"
)

func TestBaseNode_Receive(t *testing.T) {

	t.Run("Test single receive", func(t *testing.T) {
		node1 := NewBaseNode(1, 1)

		if len(node1.store) != 0 {
			t.Fatal("Store is not empty on init")
		}

		node1.addRequest(_NewRequest(_NTM("key", node1), []fp.INode{}))
		node1.Receive(_NTM("key", node1), "value")

		value, ok := node1.getFromStore("key")

		if !ok || value != "value" {
			t.Fatal("Adding to store failed")
		}
	})

	t.Run("Test receive and remember to answer", func(t *testing.T) {
		node1, node2, node3 := NewBaseNode(1, 1), NewBaseNode(1, 1), NewBaseNode(1, 1)
		node1 = node1.SetGetters(friendsFactory(node2, node3), tunnelFactory())
		node2 = node2.SetGetters(friendsFactory(node1), tunnelFactory())
		node3 = node3.SetGetters(friendsFactory(node1), tunnelFactory())

		node1.addRequest(_NewRequest(_NTM("key", node2), []fp.INode{}), _NewRequest(_NTM("key", node3), []fp.INode{}))
		node2.addRequest(_NewRequest(_NTM("key", node2), []fp.INode{node1}))
		node3.addRequest(_NewRequest(_NTM("key", node3), []fp.INode{node1}))

		node1.Receive(_NTM("key", node2), "value")

		time.Sleep(time.Millisecond * 100)

		val2, ok2 := node2.getFromStore("key")
		val3, ok3 := node3.getFromStore("key")

		if !ok2 || !ok3 || val2 != "value" || val3 != "value" {
			t.Errorf("%v %v", ok2, ok3)
			t.Fatal("Receiving failed")
		}
	})

}

func TestBaseNode_Ask(t *testing.T) {

	t.Run("Base ask", func(t *testing.T) {
		node1, node2 := NewBaseNode(1, 1), NewBaseNode(1, 1)
		node1 = node1.SetGetters(friendsFactory(node2), tunnelFactory())
		node2 = node2.SetGetters(friendsFactory(node1), tunnelFactory())

		node1.store["key"] = "value"

		node2.Ask(_NTM("key", node2))

		time.Sleep(time.Millisecond * 100)

		val, ok := node2.getFromStore("key")

		if !ok || val != "value" {
			t.Fatal("Asking failed")
		}
	})

	t.Run("Chain ask", func(t *testing.T) {
		node1, node2, node3 := NewBaseNode(1, 1), NewBaseNode(1, 1), NewBaseNode(1, 1)

		node1 = node1.SetGetters(friendsFactory(node2), tunnelFactory())
		node2 = node2.SetGetters(friendsFactory(node1, node3), tunnelFactory())
		node3 = node3.SetGetters(friendsFactory(node2), tunnelFactory())

		node3.store["key"] = "value"

		node1.Ask(_NTM("key", node1))

		time.Sleep(time.Millisecond * 20)

		val1, ok1 := node1.getFromStore("key")

		if !ok1 || val1 != "value" {
			t.Fatal("Chain ask failed")
		}
	})

	t.Run("Reject if uninterested", func(t *testing.T) {
		node1, node2 := NewBaseNode(1, 1), NewBaseNode(1, 1)

		node1 = node1.SetGetters(friendsFactory(node2), tunnelFactory()).SetWatchers(nil, func(m fp.IMessage, me fp.INode) {
			t.Fatal("Did not reject")
		})
		node2 = node2.SetGetters(friendsFactory(node1), tunnelFactory())

		node1.addRequest(_NewRequest(_NTM("key", node2), []fp.INode{}))
		node1.Receive(_NTM("key", node2), "value")

		time.Sleep(20 * time.Millisecond)
	})
}

func TestBaseNode_StopSearch(t *testing.T) {
	t.Run("Simple stop search", func(t *testing.T) {
		node1 := NewBaseNode(10, 10)
		node2 := NewBaseNode(10, 10)
		node3 := NewBaseNode(10, 10)
		node4 := NewBaseNode(10, 10)

		var messagesReceived int32 = 0
		messageUp := func(m fp.IMessage, me fp.INode) {
			atomic.StoreInt32(&messagesReceived, atomic.LoadInt32(&messagesReceived)+1)
		}

		node1 = node1.SetGetters(friendsFactory(node2, node3), func(with fp.INode) (int, int) {
			if with == node2 {
				return 10, 1
			} else {
				return 1, 1
			}
		}).SetWatchers(nil, messageUp)
		node2 = node2.SetGetters(friendsFactory(node1), func(with fp.INode) (int, int) {
			if with == node1 {
				return 10, 1
			} else {
				return 1, 1
			}
		}).SetWatchers(nil, messageUp)
		node3 = node3.SetGetters(friendsFactory(node1, node4), tunnelFactory()).SetWatchers(nil, messageUp)
		node4 = node4.SetGetters(friendsFactory(node3), tunnelFactory()).SetWatchers(nil, messageUp)

		node2.InitStore(map[string]string{"key": "value"})
		node4.InitStore(map[string]string{"key": "value"})

		node1.Ask(_NTM("key", node1))

		time.Sleep(100 * time.Millisecond)

		if messages := atomic.LoadInt32(&messagesReceived); messages != 2 {
			t.Fatalf("Expected to receive 2 message, but received - %d \n", messages)
		}
	})

	t.Run("Chain stop search", func(t *testing.T) {
		node1 := NewBaseNode(10, 10)
		node2 := NewBaseNode(10, 10)
		node3 := NewBaseNode(10, 10)
		node4 := NewBaseNode(10, 10)
		node5 := NewBaseNode(10, 10)
		node6 := NewBaseNode(10, 10)

		var messagesReceived int32 = 0
		messageUp := func(m fp.IMessage, me fp.INode) {
			t.Logf("MessageUp on %+v with %p\n", m, me)
			atomic.StoreInt32(&messagesReceived, atomic.LoadInt32(&messagesReceived)+1)
		}

		node1 = node1.SetGetters(friendsFactory(node2, node3), func(with fp.INode) (int, int) {
			if with == node2 {
				return 10, 1
			} else {
				return 1, 1
			}
		}).SetWatchers(nil, messageUp)
		node2 = node2.SetGetters(friendsFactory(node1), func(with fp.INode) (int, int) {
			if with == node1 {
				return 10, 1
			} else {
				return 1, 1
			}
		}).SetWatchers(nil, messageUp)
		node3 = node3.SetGetters(friendsFactory(node1, node4), tunnelFactory()).SetWatchers(nil, messageUp)
		node4 = node4.SetGetters(friendsFactory(node3, node5), tunnelFactory()).SetWatchers(nil, messageUp)
		node5 = node5.SetGetters(friendsFactory(node4, node6), tunnelFactory()).SetWatchers(nil, messageUp)
		node6 = node6.SetGetters(friendsFactory(node5), tunnelFactory()).SetWatchers(nil, messageUp)

		for i, n := range []fp.INode{node1, node2, node3, node4, node5, node6} {
			t.Logf("Node %d has pointer of %p\n", i+1, n)
		}

		node2.InitStore(map[string]string{"key": "valuevalue"})
		node6.InitStore(map[string]string{"key": "valuevalue"})

		node1.Ask(_NTM("key", node1))

		time.Sleep(200 * time.Millisecond)

		if messages := atomic.LoadInt32(&messagesReceived); messages != 2 {
			t.Fatalf("Expected to receive 2 message, but received - %d \n", messages)
		}
	})
}

var id int = -1

type _TESTMessage struct {
	id   int
	from fp.INode
	key  string
}

func (m _TESTMessage) Id() int {
	return m.id
}
func (m _TESTMessage) From() fp.INode {
	return m.from
}
func (m _TESTMessage) Resends() int {
	return 0
}
func (m _TESTMessage) Key() string {
	return m.key
}
func (m _TESTMessage) Done(by fp.INode) {}
func (m _TESTMessage) Resend(from fp.INode) fp.IMessage {
	m.from = from
	return m
}
func _NTM(key string, from *BaseNode) _TESTMessage {
	id++
	return _TESTMessage{key: key, from: from, id: id}
}

var friendsFactory = func(nodes ...fp.INode) func() []fp.INode {
	return func() []fp.INode {
		return nodes
	}
}
var tunnelFactory = func() func(with fp.INode) (int, int) {
	return func(with fp.INode) (int, int) {
		return 1, 1
	}
}
