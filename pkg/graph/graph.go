package graph

// Undirected unweighted graph
type Graph struct {
	grid [][]bool
}

func (g *Graph) _SetPath(n1 int, n2 int, exists bool) {
	g.grid[max(n1, n2)][min(n1, n2)] = exists
}

func (g Graph) _HasPath(n1, n2 int) bool {
	return g.grid[max(n1, n2)][min(n1, n2)]
}

func (g Graph) GetPaths(n int) []int {
	ans := []int{}

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

	return ans
}

func _NewGraph(length int) Graph {
	grid := make([][]bool, length)

	for i := 1; i < length; i++ {
		grid[i] = make([]bool, i)
	}

	return Graph{grid: grid}
}
