package graph

import "math/rand"

// gen() should not return zeroValue
func NewRandomGraph[T comparable](size int, degree int, gen func() T) Graph[T] {
	g := NewGraph[T](size)

	if size < 2 || degree >= size+1 || degree <= 0 {
		return g
	}

	probability := (float64(degree)) / float64(size-1)

	for i := 0; i < size-1; i++ {
		for j := i + 1; j < size; j++ {
			if probability >= rand.Float64() {
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
			from := cg[i][rand.Intn(len(cg[i]))]
			to := cg[i+1][rand.Intn(len(cg[i+1]))]
			g.SetPath(from, to, gen())
		}
	}

	return g
}
