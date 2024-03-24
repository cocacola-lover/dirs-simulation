package main

import (
	netp "dirs/simulation/pkg/network"
	"fmt"
	"time"
)

func main() {

	n := 10000
	net, searchers, havers := netp.NewTestSearchNetwork(n, 3, netp.SearchRequest{
		Key:               "key",
		Val:               "value",
		Popularity:        0.01,
		NumberOfSearchers: 1,
	})

	fmt.Printf("Havers are %v\n", havers)
	fmt.Printf("Searchers are %v\n", searchers)

	for i, v := range searchers[0] {
		net.Get(v).ReceiveRouteMessage(i, "key", net.Get(v))
	}

	// fmt.Print(net.Graph)

	time.Sleep(time.Millisecond * 200)

	for i, v := range searchers[0] {
		fmt.Println(net.LoggerStringById(i, net.Get(v)))
	}

}
