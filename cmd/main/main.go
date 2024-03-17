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

	net := netp.NewRandomNetwork(netp.FactoryInitRandomBaseNode(1, 10), n, 3)

	puttingMessageIn := crp.Rand.Intn(n-1) + 2
	fmt.Printf("Putting message in %d\n", puttingMessageIn)
	net.Get(puttingMessageIn).InitStore(map[string]string{"key": "value"})

	search1 := mp.NewFirstMessage("key", net.Get(0))
	search2 := mp.NewFirstMessage("key", net.Get(1))

	go net.Get(0).Ask(search1)
	go net.Get(1).Ask(search2)

	search1.WaitForDone()
	search2.WaitForDone()

	time.Sleep(time.Millisecond * 100)

	fmt.Println(net.StringLogger())
}
