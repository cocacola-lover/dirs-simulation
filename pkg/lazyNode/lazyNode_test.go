package lazynode

import (
	"dirs/simulation/pkg/node"
	"testing"
	"time"
)

func TestLazyNode(t *testing.T) {
	t.Run("Do not respond", func(t *testing.T) {
		node1 := node.NewNode(1, 1, nil, nil)
		node2 := NewLazyNode(1, 1, nil, nil)

		node1.SetOuterFunctions(
			func() []node.INode { return []node.INode{node2} },
			func(with node.INode) (int, int) { return 1, 1 },
		)
		node2.SetOuterFunctions(
			func() []node.INode { return []node.INode{node1} },
			func(with node.INode) (int, int) { return 1, 1 },
		)
		node2.PutVal("key", "value")

		node1.ReceiveRouteMessage(0, "key", node1)
		time.Sleep(50 * time.Millisecond)

		if val, ok := node1.HasKey("key"); ok {
			t.Errorf("Value is %s\n", val)
			t.Fatalf("Adding to store failed\n")
		}
	})

	t.Run("Work like usual when asking", func(t *testing.T) {
		node1 := NewLazyNode(1, 1, nil, nil)
		node2 := node.NewNode(1, 1, nil, nil)

		node1.SetOuterFunctions(
			func() []node.INode { return []node.INode{node2} },
			func(with node.INode) (int, int) { return 1, 1 },
		)
		node2.SetOuterFunctions(
			func() []node.INode { return []node.INode{node1} },
			func(with node.INode) (int, int) { return 1, 1 },
		)
		node2.PutVal("key", "value")

		node1.ReceiveRouteMessage(0, "key", node1)
		time.Sleep(50 * time.Millisecond)

		if val, ok := node1.HasKey("key"); !ok {
			t.Errorf("Value is %s\n", val)
			t.Fatalf("Adding to store failed\n")
		}
	})
}
