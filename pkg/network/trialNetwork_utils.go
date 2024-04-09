package network

import (
	crp "dirs/simulation/pkg/controlledRandom"
)

func devideSearchersAndHavers(size int, request SearchRequest) ([]int, []int) {

	hasInStore := []int{}
	usedValues := make(map[int]bool)

	for i := 0; i < size; i++ {
		if crp.Rand.Float64() <= request.Popularity {
			hasInStore = append(hasInStore, i)
			usedValues[i] = true
		}
	}

	// But at least 1 has info
	if len(hasInStore) == 0 {
		v := crp.Rand.Intn(size)
		hasInStore = append(hasInStore, v)
		usedValues[v] = true
	}

	searchers := make([]int, 0, request.NumberOfSearchers)
	for i := 0; i < request.NumberOfSearchers; i++ {
		ind := crp.Rand.Intn(size)

		for _, ok := usedValues[ind]; ok && len(hasInStore)+len(searchers) < size; _, ok = usedValues[ind] {
			ind = crp.Rand.Intn(size)
		}

		if len(hasInStore)+len(searchers) == size {
			for len(searchers) < request.NumberOfSearchers {
				searchers = append(searchers, hasInStore[len(hasInStore)-1])
				hasInStore = hasInStore[:len(hasInStore)-1]
			}
		}

		searchers = append(searchers, ind)
		usedValues[ind] = true
	}

	return searchers, hasInStore
}
