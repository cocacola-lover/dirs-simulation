package nlogger

import (
	np "dirs/simulation/pkg/node"
	"fmt"
)

func (l *Logger) StringById(id int) string {
	ans := fmt.Sprintf("For track #%d :\n", id)

	l.rmrLock.Lock()
	l.rmrdLock.Lock()
	ans += fmt.Sprintf(
		"Have %d message receives and %d receives-rejectes\n",
		countMapArr(l.routeMessageReceives[id]),
		countMapArr(l.deniedRouteMessageReceives[id]),
	)
	l.rmrdLock.Unlock()
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

	l.rmtLock.Lock()
	l.rmtdLock.Lock()
	ans += fmt.Sprintf(
		"Have %d message timeouts and %d timeouts-rejectes\n",
		countMapArr(l.routeMessageTimeouts[id]),
		countMapArr(l.deniedRouteMessageTimeouts[id]),
	)
	l.rmtdLock.Unlock()
	l.rmtLock.Unlock()

	ans += fmt.Sprintf("The download root was of length - %d\n", len(l.downloadMessages[id]))

	return ans
}

func (l *Logger) StringByIdVerbose(id int, lead np.INode, phoneBook map[np.INode]int) string {
	ans := fmt.Sprintf("For track #%d :\n", id)

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

	l.rmtLock.Lock()
	l.rmtdLock.Lock()
	ans += fmt.Sprintf(
		"Have %d message timeouts and %d timeouts-rejectes\n",
		countMapArr(l.routeMessageTimeouts[id]),
		countMapArr(l.deniedRouteMessageTimeouts[id]),
	)
	l.rmtdLock.Unlock()
	l.rmtLock.Unlock()

	l.dLock.Lock()
	ans += fmt.Sprintf("The download root was of length - %d\n", len(l.downloadMessages[id]))
	for _, el := range l.downloadMessages[id] {
		ans += fmt.Sprintf("%d -> ", phoneBook[el])
	}
	l.dLock.Unlock()

	ans += fmt.Sprint(phoneBook[lead])

	return ans
}

func (l *Logger) String() string {
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
