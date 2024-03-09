package graph

import "testing"

func TestNewRandomConnectedGraph(t *testing.T) {
	g := NewRandomConnectedGraph[bool](100, 5, func() bool { return true })

	if !g.IsConnected() {
		t.Fatal("Can not create RandomConnectedGraph")
	}
}
