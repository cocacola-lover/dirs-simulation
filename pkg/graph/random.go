package graph

import (
	crp "dirs/simulation/pkg/controlledRandom"
	"dirs/simulation/pkg/utils"
)

// gen() should not return zeroValue
func NewRandomGraph[T comparable](size int, degree int, gen func() T) Graph[T] {
	g := NewGraph[T](size)

	if size < 2 || degree >= size+1 || degree <= 0 {
		return g
	}

	probability := (float64(degree)) / float64(size-1)

	for i := 0; i < size-1; i++ {
		for j := i + 1; j < size; j++ {
			if probability >= crp.Rand.Float64() {
				g.SetPath(i, j, gen())
			}
		}
	}

	return g
}

// gen() should not return zeroValue
func NewRandomConnectedGraph[T comparable](size, degree int, gen func() T) Graph[T] {
	g := NewRandomGraph(size, degree, gen)

	cg := g.GetConnectedGroups()

	if len(cg) > 1 {
		for i := 0; i < len(cg)-1; i++ {
			from := cg[i][crp.Rand.Intn(len(cg[i]))]
			to := cg[i+1][crp.Rand.Intn(len(cg[i+1]))]
			g.SetPath(from, to, gen())
		}
	}

	return g
}

func NewScaleFreeGraph[T comparable](size, degree int, gen func() T) Graph[T] {
	g := NewGraph[T](size)

	if size < 2 || degree >= size || degree <= 0 {
		panic("Wrong input for ScaleFree Graph")
	}

	m := func(more bool) int {
		if degree%2 == 0 {
			return degree / 2
		} else if more {
			return degree/2 + 1
		} else {
			return degree / 2
		}
	}

	// Initial network
	for i := 0; i < m(true); i++ {
		for j := i + 1; j < m(true); j++ {
			g.SetPath(i, j, gen())
		}
	}

	for i := m(true); i < size; i++ {
		if i == 1 {
			g.SetPath(1, 0, gen())
			continue
		}

		distribution := make([]float64, i-1)

		var sum float64 = 0
		// Calculate sum of all degrees for each node
		for j := 0; j < i; j++ {
			sum += float64(g.GetDegree(j))
		}

		// Fill distribution
		distribution[0] = float64(g.GetDegree(0)) / sum
		for j := 1; j < i-1; j++ {
			distribution[j] = float64(g.GetDegree(j))/sum + distribution[j-1]
		}

		newLinks := make([]int, 0, m(i%2 == 0))

		// Fill new Links
		for len(newLinks) < cap(newLinks) {
			found, probability := 0, crp.Rand.Float64()

			if probability > distribution[i-2] {
				found = i - 1
			} else {
				for j := 0; j < i-1; j++ {
					if probability < distribution[j] {
						found = j
						break
					}
				}
			}

			// Check that link is not present already
			if _, has := utils.Find(newLinks, func(v int) bool {
				return found == v
			}); has {
				continue
			}

			// Otherwise add link
			newLinks = append(newLinks, found)
		}

		for _, link := range newLinks {
			g.SetPath(i, link, gen())
		}
	}

	return g
}
