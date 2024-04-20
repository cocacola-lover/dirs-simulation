package nlogger

import (
	np "dirs/simulation/pkg/node"
	"fmt"
)

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
