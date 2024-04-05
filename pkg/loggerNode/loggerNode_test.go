package loggernode

import (
	nlogger "dirs/simulation/pkg/nLogger"
	"dirs/simulation/pkg/node"
	searchernode "dirs/simulation/pkg/searcherNode"
	"log"
	"math"
	"testing"
	"time"
)

var friendsFactory = func(nodes ...node.INode) func() []node.INode {
	return func() []node.INode {
		return nodes
	}
}

var tunnelFactory = func() func(with node.INode) (int, int) {
	return func(with node.INode) (int, int) {
		return 1, 1
	}
}

func TestNode_ReceiveRouteMessage(t *testing.T) {

	t.Run("Test simple ReceiveRouteMessage", func(t *testing.T) {
		baseNode1 := node.NewNode(1, 1, nil, nil)
		baseNode2 := node.NewNode(1, 1, nil, nil)

		logger := nlogger.NewLogger()

		node1 := NewLoggerNode(searchernode.NewSearchNode(baseNode1), logger)
		node2 := NewLoggerNode(searchernode.NewSearchNode(baseNode2), logger)

		baseNode1.SetOuterFunctions(friendsFactory(node2), tunnelFactory())
		baseNode2.SetOuterFunctions(friendsFactory(node1), tunnelFactory())

		baseNode2.PutVal("key", "value")

		id := node1.StartSearchAndWatch("key")

		time.Sleep(30 * time.Millisecond)

		_, ok := baseNode1.HasKey("key")

		if !ok {
			t.Error("Adding to store failed \n")
		}

		dur, _ := logger.DurationToArriveLocked(id)
		t.Logf("Took %v to arrive", dur)

		eR, eC, eDR, eDC, eD := 2, 1, 0, 0, 1
		if logger.CountRouteMessageReceives() != eR {
			t.Errorf("Expected %d R but got %d", eR, logger.CountRouteMessageReceives())
		}
		if logger.CountRouteMessageConfirms() != eC {
			t.Errorf("Expected %d C but got %d", eC, logger.CountRouteMessageConfirms())
		}
		// if logger.CountRouteMessageTimeouts() != eT {
		// 	t.Errorf("Expected %d T but got %d", eT, logger.CountRouteMessageTimeouts())
		// }
		if logger.CountDeclinedRouteMessageReceives() != eDR {
			t.Errorf("Expected %d DR but got %d", eDR, logger.CountDeclinedRouteMessageReceives())
		}
		if logger.CountDeclinedRouteMessageConfirms() != eDC {
			t.Errorf("Expected %d DC but got %d", eDC, logger.CountDeclinedRouteMessageConfirms())
		}
		// if logger.CountDeclinedRouteMessageTimeouts() != eDT {
		// 	t.Errorf("Expected %d DT but got %d", eDT, logger.CountDeclinedRouteMessageTimeouts())
		// }
		if logger.CountDownloadMessages() != eD {
			t.Errorf("Expected %d D but got %d", eD, logger.CountDownloadMessages())
		}
	})

	t.Run("Test complex ReceiveRouteMessage", func(t *testing.T) {
		baseNode1 := node.NewNode(1, 1, nil, nil)
		baseNode2 := node.NewNode(1, 1, nil, nil)
		baseNode3 := node.NewNode(1, 1, nil, nil)
		baseNode4 := node.NewNode(1, 1, nil, nil)
		baseNode5 := node.NewNode(1, 1, nil, nil)
		baseNode6 := node.NewNode(1, 1, nil, nil)

		logger := nlogger.NewLogger()

		node1 := NewLoggerNode(searchernode.NewSearchNode(baseNode1), logger)
		node2 := NewLoggerNode(searchernode.NewSearchNode(baseNode2), logger)
		node3 := NewLoggerNode(searchernode.NewSearchNode(baseNode3), logger)
		node4 := NewLoggerNode(searchernode.NewSearchNode(baseNode4), logger)
		node5 := NewLoggerNode(searchernode.NewSearchNode(baseNode5), logger)
		node6 := NewLoggerNode(searchernode.NewSearchNode(baseNode6), logger)

		baseNode1.SetOuterFunctions(friendsFactory(node2, node3), tunnelFactory())
		baseNode2.SetOuterFunctions(friendsFactory(node1, node6), tunnelFactory())
		baseNode3.SetOuterFunctions(friendsFactory(node1, node4), tunnelFactory())
		baseNode4.SetOuterFunctions(friendsFactory(node3, node5), tunnelFactory())
		baseNode5.SetOuterFunctions(friendsFactory(node4, node6), tunnelFactory())
		baseNode6.SetOuterFunctions(friendsFactory(node5, node2), tunnelFactory())

		for i, n := range []node.INode{node1, node2, node3, node4, node5, node6} {
			log.Printf("Node%d has a pointer of %p\n", i+1, n)
		}

		baseNode4.PutVal("key", "value")

		id := node1.StartSearchAndWatch("key")

		time.Sleep(100 * time.Millisecond)

		_, ok := baseNode1.HasKey("key")

		if !ok {
			t.Error("Adding to store failed \n")
		}
		dur, _ := logger.DurationToArriveLocked(id)
		t.Logf("Took %v to arrive", dur)

		eR, eC, eDR, eDC, eD := 7, 6, 0, 0, 2
		if logger.CountRouteMessageReceives() != eR {
			t.Errorf("Expected %d R but got %d", eR, logger.CountRouteMessageReceives())
		}
		if math.Abs(float64(logger.CountRouteMessageConfirms()-eC)) > 1 {
			t.Errorf("Expected %d C but got %d", eC, logger.CountRouteMessageConfirms())
		}
		if logger.CountDeclinedRouteMessageReceives() != eDR {
			t.Errorf("Expected %d DR but got %d", eDR, logger.CountDeclinedRouteMessageReceives())
		}
		if math.Abs(float64(logger.CountDeclinedRouteMessageConfirms()-eDC)) > 1 {
			t.Errorf("Expected %d DC but got %d", eDC, logger.CountDeclinedRouteMessageConfirms())
		}
		if logger.CountDownloadMessages() != eD {
			t.Errorf("Expected %d D but got %d", eD, logger.CountDownloadMessages())
		}
	})

	t.Run("Test WaitToFinishAllSearches", func(t *testing.T) {
		baseNode1 := node.NewNode(1, 1, nil, nil)
		baseNode2 := node.NewNode(1, 1, nil, nil)

		logger := nlogger.NewLogger()

		node1 := NewLoggerNode(searchernode.NewSearchNode(baseNode1), logger)
		node2 := NewLoggerNode(searchernode.NewSearchNode(baseNode2), logger)

		baseNode1.SetOuterFunctions(friendsFactory(node2), tunnelFactory())
		baseNode2.SetOuterFunctions(friendsFactory(node1), tunnelFactory())

		baseNode2.PutVal("key", "value")

		id := node1.StartSearchAndWatch("key")

		node1.WaitToFinishAllSearches()

		_, ok := baseNode1.HasKey("key")

		if !ok {
			t.Error("Adding to store failed \n")
		}

		dur, _ := logger.DurationToArriveLocked(id)
		t.Logf("Took %v to arrive", dur)
	})
}
