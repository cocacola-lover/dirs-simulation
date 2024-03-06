package basenode

import (
	"testing"
	"time"
)

func TestBaseNode_Receive(t *testing.T) {

	t.Run("Test single receive", func(t *testing.T) {
		node1 := NewBaseNode()

		if len(node1.store) != 0 {
			t.Fatal("Store is not empty on init")
		}

		node1.Receive("key", "value")

		value, ok := node1.store["key"]

		if !ok || value != "value" {
			t.Fatal("Adding to store failed")
		}
	})

	t.Run("Test receive and remember to answer", func(t *testing.T) {
		node1 := NewBaseNode()
		node2 := NewBaseNode()
		node3 := NewBaseNode()

		node1.friends = append(node1.friends, node2)
		node2.friends = append(node2.friends, node1)
		node3.friends = append(node3.friends, node1)

		node1.requests = append(node1.requests, _Request{key: "key", from: node2}, _Request{key: "key", from: node3})

		node1.Receive("key", "value")

		time.Sleep(time.Millisecond)

		val1, ok1 := node1.store["key"]
		val2, ok2 := node2.store["key"]
		val3, ok3 := node3.store["key"]

		if !ok1 || !ok2 || !ok3 || val1 != "value" || val2 != "value" || val3 != "value" {
			t.Fatal("Receiving failed")
		}
	})

}

func TestBaseNode_Ask(t *testing.T) {
	t.Run("Base ask", func(t *testing.T) {
		node1 := NewBaseNode()
		node2 := NewBaseNode()

		node1.store["key"] = "value"

		node1.friends = append(node1.friends, node2)
		node2.friends = append(node2.friends, node1)

		node1.Ask("key", node2)

		time.Sleep(time.Millisecond)

		val, ok := node2.store["key"]

		if !ok || val != "value" {
			t.Fatal("Asking failed")
		}
	})

	t.Run("Chain ask", func(t *testing.T) {
		node1 := NewBaseNode()
		node2 := NewBaseNode()
		node3 := NewBaseNode()

		node1.friends = append(node1.friends, node2)
		node2.friends = append(node2.friends, node1, node3)
		node3.friends = append(node3.friends, node2)

		node3.store["key"] = "value"

		node2.Ask("key", node1)

		time.Sleep(time.Millisecond)

		val1, ok1 := node1.store["key"]
		val2, ok2 := node2.store["key"]

		if !ok1 || !ok2 || val1 != "value" || val2 != "value" {
			t.Fatal("Chain ask failed")
		}
	})
}
