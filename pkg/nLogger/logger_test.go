package nlogger

import (
	"dirs/simulation/pkg/node"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l := NewLogger()

	node1 := node.NewNode(0, 0, nil, nil)
	node2 := node.NewNode(0, 0, nil, nil)
	node3 := node.NewNode(0, 0, nil, nil)
	node4 := node.NewNode(0, 0, nil, nil)
	node5 := node.NewNode(0, 0, nil, nil)

	go l.AddRouteMessageReceive(0, node1, node2)
	go l.AddRouteMessageReceive(0, node1, node3)
	go l.AddRouteMessageReceive(0, node2, node4)
	go l.AddRouteMessageReceive(0, node2, node5)

	go l.AddRouteMessageReceive(1, node2, node4)
	go l.AddRouteMessageReceive(1, node2, node5)

	time.Sleep(5 * time.Millisecond)

	if l.CountRouteMessageReceives() != 6 {
		t.Fatal("Count method does not work")
	}

}
