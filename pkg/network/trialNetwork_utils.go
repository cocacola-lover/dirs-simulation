package network

import crp "dirs/simulation/pkg/controlledRandom"

func devideSearchersAndHavers(size int, request SearchRequest) ([]int, []int) {
	hasInStore := []int{}

	for i := 0; i < size; i++ {
		if crp.Rand.Float64() <= request.Popularity {
			hasInStore = append(hasInStore, i)
		}
	}

	// But at least 1 has info
	if len(hasInStore) == 0 {
		hasInStore = append(hasInStore, crp.Rand.Intn(size))
	}

	searchers := []int{}
	for i := 0; i < request.NumberOfSearchers; i++ {
		ind := crp.Rand.Intn(size - len(hasInStore) - len(searchers))

		jumpOver := 0
		for h, s := 0, 0; h+s < len(hasInStore)+len(searchers); {
			var v int
			if len(searchers) > s && hasInStore[h] > searchers[s] {
				v = searchers[s]
				s++
			} else {
				v = hasInStore[h]
				h++
			}

			if v > ind+jumpOver {
				break
			}
			jumpOver++
		}

		searchers = append(searchers, ind+jumpOver)
	}

	return searchers, hasInStore
}
