package graph

import "math/rand"

func NewRandomGraph(size, degree int) Graph {
	g := NewGraph(size)

	if size < 2 || degree >= size+1 || degree <= 0 {
		return g
	}

	probability := (float64(degree)) / float64(size-1)

	for i := 0; i < size-1; i++ {
		for j := i + 1; j < size; j++ {
			g.SetPath(i, j, probability >= rand.Float64())
		}
	}

	return g
}

func NewRandomConnectedGraph(size, degree int) Graph {
	g := NewRandomGraph(size, degree)

	cg := g.GetConnectedGroups()

	if len(cg) > 1 {
		for i := 0; i < len(cg)-1; i++ {
			from := cg[i][rand.Intn(len(cg[i]))]
			to := cg[i+1][rand.Intn(len(cg[i+1]))]
			g.SetPath(from, to, true)
		}
	}

	return g
}
