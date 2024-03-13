package basenode

import (
	"dirs/simulation/pkg/network"
	netp "dirs/simulation/pkg/network"
	"testing"
	"time"
)

func TestBaseNode_Receive(t *testing.T) {

	t.Run("Test single receive", func(t *testing.T) {
		net := netp.NewEmptyNetwork(func(net *network.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 1)

		if len(net.Get(0).store) != 0 {
			t.Fatal("Store is not empty on init")
		}

		net.Get(0).Receive("key", "value")

		value, ok := net.Get(0).store["key"]

		if !ok || value != "value" {
			t.Fatal("Adding to store failed")
		}
	})

	t.Run("Test receive and remember to answer", func(t *testing.T) {

		net := netp.NewEmptyNetwork(func(net *network.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 3)

		net.SetPath(0, 1, netp.Tunnel{Length: 1, Width: 1})
		net.SetPath(0, 2, netp.Tunnel{Length: 1, Width: 1})

		net.Get(0).requests = append(net.Get(0).requests, _Request{key: "key", from: net.Get(1)}, _Request{key: "key", from: net.Get(2)})

		net.Get(0).Receive("key", "value")

		time.Sleep(time.Millisecond * 100)

		val1, ok1 := net.Get(0).store["key"]
		val2, ok2 := net.Get(1).store["key"]
		val3, ok3 := net.Get(2).store["key"]

		if !ok1 || !ok2 || !ok3 || val1 != "value" || val2 != "value" || val3 != "value" {
			t.Errorf("%v %v %v", ok1, ok2, ok3)
			t.Fatal("Receiving failed")
		}
	})

}

func TestBaseNode_Ask(t *testing.T) {
	t.Run("Base ask", func(t *testing.T) {
		net := netp.NewEmptyNetwork(func(net *network.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 2)

		net.Get(0).store["key"] = "value"

		net.SetPath(0, 1, netp.Tunnel{Length: 1, Width: 1})

		net.Get(0).Ask("key", net.Get(1))

		time.Sleep(time.Millisecond * 10)

		val, ok := net.Get(1).store["key"]

		if !ok || val != "value" {
			t.Fatal("Asking failed")
		}
	})

	t.Run("Chain ask", func(t *testing.T) {
		net := netp.NewEmptyNetwork(func(net *network.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 3)

		net.SetPath(0, 1, netp.Tunnel{Length: 1, Width: 1})
		net.SetPath(1, 2, netp.Tunnel{Length: 1, Width: 1})

		net.Get(2).store["key"] = "value"

		net.Get(1).Ask("key", net.Get(0))

		time.Sleep(time.Millisecond * 12)

		val1, ok1 := net.Get(0).store["key"]
		val2, ok2 := net.Get(1).store["key"]

		if !ok1 || !ok2 || val1 != "value" || val2 != "value" {
			t.Fatal("Chain ask failed")
		}
	})
}
