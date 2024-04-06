package main

import (
	nlogger "dirs/simulation/pkg/nLogger"
	netp "dirs/simulation/pkg/network"
	"fmt"
	"time"
)

func runExperiment(logger *nlogger.Logger) {
	n := 400

	trialNet := netp.NewTrialNetwork(netp.NewBaseNetwork(n, 3, logger))

	trialNet.RunRequests([]netp.SearchRequest{{
		Id:                0,
		Key:               fmt.Sprint(0),
		Val:               "val",
		Popularity:        0.05,
		NumberOfSearchers: 8,
	}})
	time.Sleep(time.Millisecond * 10)

	trialNet.WaitToFinishAllSearchers()
	trialNet.Close()
}

func main() {

	// n := 1000

	logger := nlogger.NewLogger()

	for i := 0; i < 20; i++ {
		runExperiment(logger)
	}

	// trialNet := netp.NewTrialNetwork(netp.NewBaseNetwork(n, 4, logger))

	// trialNet.RunRequests([]netp.SearchRequest{{
	// 	Id:                0,
	// 	Key:               fmt.Sprint(0),
	// 	Val:               "val",
	// 	Popularity:        0.01,
	// 	NumberOfSearchers: n / 20,
	// }})
	// time.Sleep(time.Millisecond * 10)

	// trialNet.WaitToFinishAllSearchers()

	fmt.Print(logger.String())
}
