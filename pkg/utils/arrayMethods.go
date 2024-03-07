package utils

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
