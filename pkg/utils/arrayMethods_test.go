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
