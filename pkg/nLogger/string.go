package nlogger

import (
	np "dirs/simulation/pkg/node"
	"fmt"
	"sync/atomic"
)

func (l *Logger) StringById(id int) string {
	ans := fmt.Sprintf("For track #%d :\n", id)

	ans += fmt.Sprintf(
		"Have %d message receives and %d receives-rejectes\n",
		countMapArrWithLock(l.routeMessageReceives[id], &l.rmrLock),
		countMapArrWithLock(l.deniedRouteMessageReceives[id], &l.rmrdLock),
	)

	ans += fmt.Sprintf(
		"Have %d message confirms and %d confirms-rejectes\n",
		countMapArrWithLock(l.routeMessageConfirms[id], &l.rmcLock),
		countMapArrWithLock(l.deniedRouteMessageConfirms[id], &l.rmcdLock),
	)

	ans += fmt.Sprintf(
		"The download root was of length - %d, but had %d downloads\n",
		l.CountDownloadPath(id),
		len(l.downloadMessages[id]),
	)

	dur, ok := l.DurationToArriveLocked(id)
	if ok {
		ans += fmt.Sprintf("The download took %v\n", dur)
	} else {
		ans += "WARNING : message never reached root"
	}

	return ans
}

func (l *Logger) StringByIdVerbose(id int, phoneBook map[np.INode]int) string {
	ans := fmt.Sprintf("For track #%d :\n", id)

	l.dLock.Lock()
	lead, ok := l.startedSearches[id]
	if !ok {
		panic("Lead for search does not exist")
	}
	l.dLock.Unlock()

	l.rmrLock.Lock()
	l.rmrdLock.Lock()
	ans += fmt.Sprintf(
		"Have %d message receives and %d receives-rejectes\n",
		countMapArr(l.routeMessageReceives[id]),
		countMapArr(l.deniedRouteMessageReceives[id]),
	)
	l.rmrdLock.Unlock()

	nodes := orderFromLead(l.routeMessageReceives[id], lead)
	for i, narr := range nodes {
		for _, el := range narr {
			ans += fmt.Sprintf("%d ", phoneBook[el])
		}
		if i != len(nodes)-1 {
			ans += "-> "
		} else {
			ans += "\n"
		}
	}
	l.rmrLock.Unlock()

	l.rmcLock.Lock()
	l.rmcdLock.Lock()
	ans += fmt.Sprintf(
		"Have %d message confirms and %d confirms-rejectes\n",
		countMapArr(l.routeMessageConfirms[id]),
		countMapArr(l.deniedRouteMessageConfirms[id]),
	)
	l.rmcdLock.Unlock()
	l.rmcLock.Unlock()

	ans += fmt.Sprintf(
		"Have %d fault message receives\n",
		atomic.LoadInt32(&l.faultMessageReceives),
	)

	l.dLock.Lock()
	ans += fmt.Sprintf("The download root was of length - %d\n", len(l.downloadMessages[id]))
	for _, el := range l.downloadMessages[id] {
		ans += fmt.Sprintf("%d -> ", phoneBook[el])
	}
	l.dLock.Unlock()

	ans += fmt.Sprintf("%v\n", phoneBook[lead])

	dur, ok := l.DurationToArriveLocked(id)
	if ok {
		ans += fmt.Sprintf("The download took %v\n", dur)
	} else {
		ans += "WARNING!!!!!!!!!!!!!!!!!!!!!!!!!!!!! : message never reached root\n\n\n"
	}

	return ans
}

func (l *Logger) StringByIdForEach() string {
	ans := ""

	l.dLock.Lock()
	keys := []int{}
	for key := range l.downloadMessages {
		keys = append(keys, key)
	}
	l.dLock.Unlock()

	for key := range keys {
		ans += l.StringById(key)
	}

	return ans
}

func (l *Logger) StringByIdForEachVerbose(phoneBook map[np.INode]int) string {
	ans := ""

	l.dLock.Lock()
	keys := []int{}
	for key := range l.downloadMessages {
		keys = append(keys, key)
	}
	l.dLock.Unlock()

	for key := range keys {
		ans += l.StringByIdVerbose(key, phoneBook)
	}

	return ans
}

func (l *Logger) String() string {

	ans := "Summary of the experiment :\n"

	averageMessages := l.AverageRouteMessageReceives() + l.AverageDeclinedRouteMessageReceives()
	averageMessages += l.AverageRouteMessageConfirms() + l.AverageDeclinedRouteMessageConfirms()
	averageMessages += l.AverageFaultMessageReceives() + l.AverageDownloadMessages()

	ans += fmt.Sprintf("Average messages sent - %v\n", averageMessages)

	ans += fmt.Sprintf(
		"The average download root was of length - %v\n",
		l.AverageDownloadPath(),
	)

	dur, _ := l.AverageDurationToArriveLocked()
	ans += fmt.Sprintf("The average download took %v\n", dur)

	return ans
}

// Private ----------------------------------------------------------------------------

func orderFromLead(dict map[np.INode][]np.INode, lead np.INode) [][]np.INode {
	ans := [][]np.INode{{lead}}

	for i := 0; true; i++ {
		newArr := []np.INode{}

		for _, n := range ans[i] {
			if dictArr, ok := dict[n]; ok {
				for _, el := range dictArr {
					if el == n {
						continue
					}
					newArr = append(newArr, el)
				}
			}
		}

		if len(newArr) == 0 {
			break
		} else {
			ans = append(ans, newArr)
		}
	}

	return ans
}
