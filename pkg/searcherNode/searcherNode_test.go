package searchernode

import (
	"dirs/simulation/pkg/node"
	"testing"
)

func TestSearchNode(t *testing.T) {
	// bn := node.NewNode()

	t.Run("Test simple ReceiveRouteMessage", func(t *testing.T) {
		node1 := node.NewNode(1, 1, nil, nil)
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

		searchNode1 := NewSearchNode(node1)
		stopCh := make(chan bool)
		searchNode1.StartSearch(0, "key", stopCh)

		_, ok1 := <-stopCh
		value, ok := node1.HasKey("key")

		if ok1 || !ok || value != "value" {
			t.Fatalf("Adding to store failed\n")
		}
	})
}
