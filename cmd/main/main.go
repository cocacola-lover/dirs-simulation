package main

import (
	crp "dirs/simulation/pkg/controlledRandom"
	netp "dirs/simulation/pkg/network"
	"fmt"
	"strings"
	"time"
)

func main() {

	id := 0
	n := 10000

	trialNet := netp.NewTrialNetwork(netp.NewBaseNetwork(n, 3))

	trialNet.GenerateTasks(func() netp.SearchRequest {
		defer func() {
			id++
		}()

		return netp.SearchRequest{
			Id:                id,
			Key:               fmt.Sprint(id),
			Val:               strings.Repeat("v", 4*(id%4)),
			Popularity:        crp.Rand.Float64() / 20,
			NumberOfSearchers: crp.Rand.Intn(9) + 1,
		}
	}, func() time.Duration {
		return time.Millisecond * 100
	}, time.Second)

	time.Sleep(time.Millisecond * 100)

	fmt.Print(trialNet)
}
