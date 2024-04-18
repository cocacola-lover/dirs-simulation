package main

import (
	nlogger "dirs/simulation/pkg/nLogger"
	netp "dirs/simulation/pkg/network"
	"fmt"
	"time"
)

func runExperiment(logger *nlogger.Logger) int {
	n := 1000

	trialNet := netp.NewTrialNetwork(netp.NewFailingNetwork(n, 3, 0.01, logger))

	// for i := 0; i < n; i++ {
	// 	fmt.Printf("Node%d has a pointer of %p\n", i, trialNet.Get(i))
	// }

	havers, _ := trialNet.RunRequest(netp.SearchRequest{
		Key:               fmt.Sprint(0),
		Val:               "val",
		Popularity:        0.005,
		NumberOfSearchers: 100,
	})

	time.Sleep(100 * time.Millisecond)

	select {
	case <-time.After(time.Second * 2):
	case <-trialNet.WaitToFinishAllSearchers():
	}

	failedHavers := 0

	for _, i := range havers {
		if trialNet.Get(i).HasFailed() {
			failedHavers++
		}
	}

	trialNet.Close()

	return failedHavers
}

func main() {

	expN := 10

	logger := nlogger.NewLogger()

	var failedHavers float64 = 0
	for i := 0; i < expN; i++ {
		failedHavers += float64(runExperiment(logger))
	}

	fmt.Print(logger.String())

	failed := logger.CountFailedNodes()
	fmt.Printf("On average %v node failed\n", float64(failed)/float64(expN))
	fmt.Printf("On average %v havers failed\n", failedHavers/float64(expN))

	_, didntReach := logger.AverageDurationToArriveLocked()
	if didntReach != 0 {
		fmt.Printf("WARNING : on average %v message never reached root\n", float64(didntReach)/float64(expN))
	}
}
