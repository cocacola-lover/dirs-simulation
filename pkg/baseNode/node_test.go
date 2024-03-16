package basenode

import (
	mp "dirs/simulation/pkg/message"
	netp "dirs/simulation/pkg/network"
	"testing"
	"time"
)

func TestBaseNode_Receive(t *testing.T) {

	t.Run("Test single receive", func(t *testing.T) {
		net := netp.NewEmptyNetwork(func(net *netp.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 1)

		if len(net.Get(0).store) != 0 {
			t.Fatal("Store is not empty on init")
		}

		net.Get(0).AddRequest(mp.NewBaseMessage("key", net.Get(0)))
		net.Get(0).Receive(mp.NewBaseMessage("key", net.Get(0)), "value")

		value, ok := net.Get(0).GetFromStore("key")

		if !ok || value != "value" {
			t.Fatal("Adding to store failed")
		}
	})

	t.Run("Test receive and remember to answer", func(t *testing.T) {

		net := netp.NewEmptyNetwork(func(net *netp.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 3)

		net.SetPath(0, 1, netp.Tunnel{Length: 1, Width: 1})
		net.SetPath(0, 2, netp.Tunnel{Length: 1, Width: 1})

		net.Get(0).AddRequest(mp.NewBaseMessage("key", net.Get(1)), mp.NewBaseMessage("key", net.Get(2)))
		net.Get(1).AddRequest(mp.NewBaseMessage("key", net.Get(1)))
		net.Get(2).AddRequest(mp.NewBaseMessage("key", net.Get(2)))

		net.Get(0).Receive(mp.NewBaseMessage("key", net.Get(1)), "value")

		time.Sleep(time.Millisecond * 100)

		val2, ok2 := net.Get(1).GetFromStore("key")
		val3, ok3 := net.Get(2).GetFromStore("key")

		if !ok2 || !ok3 || val2 != "value" || val3 != "value" {
			t.Errorf("%v %v", ok2, ok3)
			t.Fatal("Receiving failed")
		}
	})

}

func TestBaseNode_Ask(t *testing.T) {
	t.Run("Base ask", func(t *testing.T) {
		net := netp.NewEmptyNetwork(func(net *netp.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 2)

		net.Get(0).store["key"] = "value"

		net.SetPath(0, 1, netp.Tunnel{Length: 1, Width: 1})

		net.Get(1).Ask(mp.NewBaseMessage("key", net.Get(1)))

		time.Sleep(time.Millisecond * 100)

		val, ok := net.Get(1).GetFromStore("key")

		if !ok || val != "value" {
			t.Fatal("Asking failed")
		}
	})

	t.Run("Chain ask", func(t *testing.T) {
		net := netp.NewEmptyNetwork(func(net *netp.Network[BaseNode], i int) *BaseNode {
			return NewBaseNode(net, 1, 1)
		}, 3)

		net.SetPath(0, 1, netp.Tunnel{Length: 1, Width: 1})
		net.SetPath(1, 2, netp.Tunnel{Length: 1, Width: 1})

		net.Get(2).store["key"] = "value"

		net.Get(0).Ask(mp.NewBaseMessage("key", net.Get(0)))

		time.Sleep(time.Millisecond * 14)

		val1, ok1 := net.Get(0).GetFromStore("key")

		if !ok1 || val1 != "value" {
			t.Fatal("Chain ask failed")
		}
	})
}
