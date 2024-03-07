package graph

import (
	"fmt"
	"slices"
)

// Undirected unweighted graph
type Graph struct {
	grid [][]bool
}

func (g Graph) Len() int {
	return len(g.grid)
}

func (g Graph) _IsOutOfBounds(args ...int) bool {
	return slices.Max(args) >= g.Len()
}

// Returns ok - false, if index out of bounds
func (g *Graph) _SetPath(n1 int, n2 int, exists bool) bool {
	if g._IsOutOfBounds(n1, n2) {
		return false
	}
	g.grid[max(n1, n2)][min(n1, n2)] = exists
	return true
}

// Returns ok - false, if index out of bounds
func (g Graph) HasPath(n1, n2 int) (bool, bool) {
	if g._IsOutOfBounds(n1, n2) {
		return false, false
	}
	return g.grid[max(n1, n2)][min(n1, n2)], true
}

func (g Graph) IsConnected() bool {
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
			if g.grid[n][i] {
				recursiveSearch(i)
			}
		}

		for i := n + 1; i < len(g.grid); i++ {
			if g.grid[i][n] {
				recursiveSearch(i)
			}
		}
	}

	recursiveSearch(0)

	return len(set) == len(g.grid)
}

func (g Graph) GetConnectedGroups() [][]int {
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
				if g.grid[n][i] {
					recursiveSearch(i)
				}
			}

			for i := n + 1; i < len(g.grid); i++ {
				if g.grid[i][n] {
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
func (g Graph) GetPaths(n int) ([]int, bool) {
	ans := []int{}
	if g._IsOutOfBounds(n) {
		return ans, false
	}

	for i := 0; i < n; i++ {
		if g.grid[n][i] {
			ans = append(ans, i)
		}
	}

	for i := n + 1; i < len(g.grid); i++ {
		if g.grid[i][n] {
			ans = append(ans, i)
		}
	}

	return ans, true
}

func NewGraph(size int) Graph {
	grid := make([][]bool, size)

	for i := 1; i < size; i++ {
		grid[i] = make([]bool, i)
	}

	return Graph{grid: grid}
}

func (g Graph) String() string {
	ans := "\n"

	for i := 0; i < len(g.grid); i++ {
		paths, _ := g.GetPaths(i)
		ans += fmt.Sprintf("%d : %v \n", i, paths)
	}

	return ans
}
