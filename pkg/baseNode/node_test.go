package basenode

import (
	"testing"
	"time"
)

func TestBaseNode_Receive(t *testing.T) {

	t.Run("Test single receive", func(t *testing.T) {
		node1 := NewBaseNode(1, 1)

		if len(node1.store) != 0 {
			t.Fatal("Store is not empty on init")
		}

		node1.AddRequest(_NewTestMessage("key", node1))
		node1.Receive(_NewTestMessage("key", node1), "value")

		value, ok := node1.GetFromStore("key")

		if !ok || value != "value" {
			t.Fatal("Adding to store failed")
		}
	})

	t.Run("Test receive and remember to answer", func(t *testing.T) {
		node1, node2, node3 := NewBaseNode(1, 1), NewBaseNode(1, 1), NewBaseNode(1, 1)

		getTunnel := func(with *BaseNode) (int, int) {
			return 1, 1
		}

		node1, node2, node3 = node1.SetGetters(func() []*BaseNode {
			return []*BaseNode{node2, node3}
		}, getTunnel), node2.SetGetters(func() []*BaseNode {
			return []*BaseNode{node1}
		}, getTunnel), node3.SetGetters(func() []*BaseNode {
			return []*BaseNode{node1}
		}, getTunnel)

		node1.AddRequest(_NewTestMessage("key", node2), _NewTestMessage("key", node3))
		node2.AddRequest(_NewTestMessage("key", node2))
		node3.AddRequest(_NewTestMessage("key", node3))

		node1.Receive(_NewTestMessage("key", node2), "value")

		time.Sleep(time.Millisecond * 100)

		val2, ok2 := node2.GetFromStore("key")
		val3, ok3 := node3.GetFromStore("key")

		if !ok2 || !ok3 || val2 != "value" || val3 != "value" {
			t.Errorf("%v %v", ok2, ok3)
			t.Fatal("Receiving failed")
		}
	})

}

func TestBaseNode_Ask(t *testing.T) {
	t.Run("Base ask", func(t *testing.T) {
		node1, node2 := NewBaseNode(1, 1), NewBaseNode(1, 1)
		node1, node2 = node1.SetGetters(func() []*BaseNode {
			return []*BaseNode{node2}
		}, func(with *BaseNode) (int, int) { return 1, 1 }), node2.SetGetters(func() []*BaseNode {
			return []*BaseNode{node1}
		}, func(with *BaseNode) (int, int) { return 1, 1 })

		node1.store["key"] = "value"

		node2.Ask(_NewTestMessage("key", node2))

		time.Sleep(time.Millisecond * 100)

		val, ok := node2.GetFromStore("key")

		if !ok || val != "value" {
			t.Fatal("Asking failed")
		}
	})

	t.Run("Chain ask", func(t *testing.T) {
		node1, node2, node3 := NewBaseNode(1, 1), NewBaseNode(1, 1), NewBaseNode(1, 1)
		node1, node2, node3 = node1.SetGetters(func() []*BaseNode {
			return []*BaseNode{node2}
		}, func(with *BaseNode) (int, int) { return 1, 1 }), node2.SetGetters(func() []*BaseNode {
			return []*BaseNode{node1, node3}
		}, func(with *BaseNode) (int, int) { return 1, 1 }), node3.SetGetters(func() []*BaseNode {
			return []*BaseNode{node2}
		}, func(with *BaseNode) (int, int) { return 1, 1 })

		node3.store["key"] = "value"

		node1.Ask(_NewTestMessage("key", node1))

		time.Sleep(time.Millisecond * 14)

		val1, ok1 := node1.GetFromStore("key")

		if !ok1 || val1 != "value" {
			t.Fatal("Chain ask failed")
		}
	})
}

var id int = -1

type _TESTMessage struct {
	id   int
	from *BaseNode
	key  string
}

func (m _TESTMessage) Id() int {
	return m.id
}
func (m _TESTMessage) From() *BaseNode {
	return m.from
}
func (m _TESTMessage) Key() string {
	return m.key
}
func (m _TESTMessage) IsValid() bool {
	return true
}
func (m _TESTMessage) Resend(from *BaseNode) any {
	m.from = from
	return m
}
func _NewTestMessage(key string, from *BaseNode) _TESTMessage {
	id++
	return _TESTMessage{key: key, from: from, id: id}
}
