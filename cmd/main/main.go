package main

import (
	crp "dirs/simulation/pkg/controlledRandom"
	mp "dirs/simulation/pkg/message"
	netp "dirs/simulation/pkg/network"
	"fmt"
	"time"
)

func main() {

	n := 100

	net := netp.NewRandomNetwork(netp.FactoryInitRandomBaseNode(1, 10), n, 5)

	puttingMessageIn := crp.Rand.Intn(n-1) + 1

	fmt.Printf("Putting message in %d\n", puttingMessageIn)

	net.Get(puttingMessageIn).InitStore(map[string]string{"key": "value"})

	net.Get(0).Ask(mp.NewBaseMessage("key", net.Get(0)))

	time.Sleep(time.Millisecond * 200)

	fmt.Println(net.StringLogger())
}
