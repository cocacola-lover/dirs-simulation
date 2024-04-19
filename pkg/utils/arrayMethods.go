package utils

import "math/rand"

func Find[T any](arr []T, fu func(v T) bool) (T, bool) {

	var zeroValue T

	if len(arr) == 0 {
		return zeroValue, false
	}

	for _, v := range arr {
		if fu(v) {
			return v, true
		}
	}
	return zeroValue, false
}

func SameValues[T comparable](arr1 []T, arr2 []T) bool {
	set := make(map[T]bool, min(len(arr1), len(arr2)))

	for _, v := range arr1 {
		set[v] = false
	}

	for _, v := range arr2 {
		_, isHere := set[v]
		if !isHere {
			return false
		}

		set[v] = true
	}

	for _, v := range set {
		if !v {
			return false
		}
	}

	return true
}

// Filters out elements for which filter fn returns false.
// Filter does not keep the order of the elements.
func Filter[T any](arr []T, filterFn func(val T, ind int) bool) []T {
	ans := []T{}

	for i, val := range arr {
		if filterFn(val, i) {
			ans = append(ans, val)
		}
	}

	return ans
}

// Shuffles in place
func Shuffle[T any](arr []T) []T {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}
