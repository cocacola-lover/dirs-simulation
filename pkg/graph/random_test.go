package graph

import "testing"

func TestNewRandomConnectedGraph(t *testing.T) {
	g := NewRandomConnectedGraph(100, 5)

	if !g.IsConnected() {
		t.Fatal("Can not create RandomConnectedGraph")
	}
}
