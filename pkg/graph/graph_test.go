package graph

import (
	"dirs/simulation/pkg/utils"
	"testing"
)

func TestGraph_GetPaths(t *testing.T) {
	g := _NewGraph(5)

	g._SetPath(3, 2, true)
	g._SetPath(3, 4, true)
	g._SetPath(3, 1, true)

	test1 := g.GetPaths(1)
	test2 := g.GetPaths(2)
	test3 := g.GetPaths(3)
	test4 := g.GetPaths(4)

	ans1 := []int{3}
	ans2 := []int{3}
	ans3 := []int{1, 2, 4}
	ans4 := []int{3}

	if !utils.SameValues(test1, ans1) {
		t.Fatalf("First test failed. Expected : %v - but got : %v", test1, ans1)
	}

	if !utils.SameValues(test2, ans2) {
		t.Fatalf("Second test failed. Expected : %v - but got : %v", test2, ans2)
	}

	if !utils.SameValues(test3, ans3) {
		t.Fatalf("Third test failed. Expected : %v - but got : %v", test3, ans3)
	}

	if !utils.SameValues(test4, ans4) {
		t.Fatalf("Forth test failed. Expected : %v - but got : %v", test4, ans4)
	}
}

func Test_IsConnected(t *testing.T) {
	g := _NewGraph(6)

	if g.IsConnected() {
		t.Fatal("Error on test 1")
	}

	g._SetPath(3, 2, true)
	g._SetPath(3, 4, true)
	g._SetPath(3, 1, true)

	if g.IsConnected() {
		t.Fatal("Error on test 2")
	}

	g._SetPath(0, 5, true)
	g._SetPath(2, 0, true)

	if !g.IsConnected() {
		t.Fatal("Error on test 3")
	}

}
