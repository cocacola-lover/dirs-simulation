package utils

import "sync"

func WithLocked[ReturnType any](lock *sync.Mutex, fu func() ReturnType) ReturnType {
	lock.Lock()
	defer lock.Unlock()
	return fu()
}

func WithLockedNoResult(lock *sync.Mutex, fu func()) {
	lock.Lock()
	defer lock.Unlock()
	fu()
}

func ZeroValue[T any]() T {
	var ans T
	return ans
}
