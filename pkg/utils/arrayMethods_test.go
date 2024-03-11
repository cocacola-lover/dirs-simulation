package utils

import "testing"

func TestSameValues(t *testing.T) {
	if !SameValues([]int{1, 2, 3}, []int{3, 1, 2}) {
		t.Fatal("Error on test 1")
	}

	if SameValues([]int{1, 2, 3, 4}, []int{3, 1, 2}) {
		t.Fatal("Error on test 2")
	}

	if !SameValues([]int{1, 2, 3, 3}, []int{3, 1, 2}) {
		t.Fatal("Error on test 3")
	}
}

func TestFilter(t *testing.T) {
	if !SameValues(Filter([]int{1, 2, 3, 4, 5, 6}, func(val, ind int) bool {
		return val%2 == 0
	}), []int{2, 4, 6}) {
		t.Fatalf("Error on test 1. Got : %v ; But Expected : %v", Filter([]int{1, 2, 3, 4, 5, 6}, func(val, ind int) bool {
			return val%2 == 0
		}), []int{2, 4, 6})
	}
}
