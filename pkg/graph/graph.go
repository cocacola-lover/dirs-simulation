package graph

import (
	"dirs/simulation/pkg/utils"
	"fmt"
	"slices"
)

// Undirected unweighted graph
type Graph[T comparable] struct {
	grid [][]T
}

func (g Graph[T]) Len() int {
	return len(g.grid)
}

func (g Graph[T]) _IsOutOfBounds(args ...int) bool {
	return slices.Max(args) >= g.Len()
}

// Returns ok - false, if index out of bounds
func (g *Graph[T]) SetPath(n1 int, n2 int, val T) bool {
	if g._IsOutOfBounds(n1, n2) {
		return false
	}
	g.grid[max(n1, n2)][min(n1, n2)] = val
	return true
}

// Returns ok - false, if index out of bounds
func (g Graph[T]) HasPath(n1, n2 int) (T, bool) {
	if g._IsOutOfBounds(n1, n2) {
		return utils.ZeroValue[T](), false
	}
	return g.grid[max(n1, n2)][min(n1, n2)], true
}

func (g Graph[T]) IsConnected() bool {
	if len(g.grid) == 0 {
		return true
	}

	set := make(map[int]bool, len(g.grid))

	var recursiveSearch func(n int)
	recursiveSearch = func(n int) {
		_, isHere := set[n]
		if isHere {
			return
		}

		set[n] = true

		for i := 0; i < n; i++ {
			if g.grid[n][i] != utils.ZeroValue[T]() {
				recursiveSearch(i)
			}
		}

		for i := n + 1; i < len(g.grid); i++ {
			if g.grid[i][n] != utils.ZeroValue[T]() {
				recursiveSearch(i)
			}
		}
	}

	recursiveSearch(0)

	return len(set) == len(g.grid)
}

func (g Graph[T]) GetConnectedGroups() [][]int {
	if len(g.grid) == 0 {
		return make([][]int, 1)
	}

	ans := [][]int{}
	set := make(map[int]bool, len(g.grid))

	// While not all nodes found
	for len(set) != len(g.grid) {
		// Find unconnected target
		var target int
		for i := 0; i < len(g.grid); i++ {
			if _, isHere := set[i]; !isHere {
				target = i
				break
			}
		}

		group := []int{}

		var recursiveSearch func(n int)
		recursiveSearch = func(n int) {
			_, isHere := set[n]
			if isHere {
				return
			}

			set[n] = true
			group = append(group, n)

			for i := 0; i < n; i++ {
				if g.grid[n][i] != utils.ZeroValue[T]() {
					recursiveSearch(i)
				}
			}

			for i := n + 1; i < len(g.grid); i++ {
				if g.grid[i][n] != utils.ZeroValue[T]() {
					recursiveSearch(i)
				}
			}
		}

		recursiveSearch(target)
		ans = append(ans, group)
	}

	return ans
}

// Returns ok - false, if index out of bounds
func (g Graph[T]) GetPaths(n int) ([]int, bool) {
	ans := []int{}
	if g._IsOutOfBounds(n) {
		return ans, false
	}

	for i := 0; i < n; i++ {
		if g.grid[n][i] != utils.ZeroValue[T]() {
			ans = append(ans, i)
		}
	}

	for i := n + 1; i < len(g.grid); i++ {
		if g.grid[i][n] != utils.ZeroValue[T]() {
			ans = append(ans, i)
		}
	}

	return ans, true
}

// Returns ok - false, if index out of bounds
func (g Graph[T]) GetDegree(n int) int {
	ans := 0
	if g._IsOutOfBounds(n) {
		panic("GetDegree index out of bounds")
	}

	for i := 0; i < n; i++ {
		if g.grid[n][i] != utils.ZeroValue[T]() {
			ans++
		}
	}

	for i := n + 1; i < len(g.grid); i++ {
		if g.grid[i][n] != utils.ZeroValue[T]() {
			ans++
		}
	}

	return ans
}

func NewGraph[T comparable](size int) Graph[T] {
	grid := make([][]T, size)

	for i := 1; i < size; i++ {
		grid[i] = make([]T, i)
	}

	return Graph[T]{grid: grid}
}

func (g Graph[T]) String() string {
	ans := "\n"

	for i := 0; i < len(g.grid); i++ {
		paths, _ := g.GetPaths(i)
		ans += fmt.Sprintf("%d : %v \n", i, paths)
	}

	return ans
}
