package message

import (
	"sync"
)

var id int = -1
var mu sync.Mutex

type BaseMessage[T any] struct {
	Id      int
	Key     string
	ReSends int
	From    *T
}

func (m BaseMessage[T]) IsValid() bool {
	return m.ReSends <= 10
}

func (m BaseMessage[T]) Resend(from *T) BaseMessage[T] {
	m.ReSends++
	m.From = from

	return m
}

func NewBaseMessage[T any](key string, from *T) BaseMessage[T] {
	mu.Lock()
	defer mu.Unlock()

	id++

	ans := BaseMessage[T]{Id: id, Key: key, ReSends: -1, From: from}

	return ans
}
