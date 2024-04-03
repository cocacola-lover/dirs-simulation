package main

import (
	netp "dirs/simulation/pkg/network"
	"fmt"
	"time"
)

func main() {

	id := 0
	n := 100

	trialNet := netp.NewTrialNetwork(netp.NewBaseNetwork(n, 4))

	trialNet.RunRequests([]netp.SearchRequest{{
		Id:                id,
		Key:               fmt.Sprint(id),
		Val:               "value",
		Popularity:        0.01,
		NumberOfSearchers: 5,
	}})

	time.Sleep(time.Second * 2)

	fmt.Print(trialNet)
}
