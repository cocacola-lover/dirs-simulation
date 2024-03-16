package logger

import (
	mp "dirs/simulation/pkg/message"
	"fmt"
	"sync"
	"sync/atomic"
)

// # The purpose of the logger is to keep track of operations for statistic's sake.
//
// ## The logger should count the number of messages sent between nodes.
// ## The logger should keep track of searches in the network and their results.
type Logger[T any] struct {
	messagesSent int32

	messageTracks map[int][]_Record[T]
	mtMu          sync.RWMutex
}

func (l *Logger[T]) MessagesSent() int {
	return int(atomic.LoadInt32(&l.messagesSent))
}

func (l *Logger[T]) AddMessage(m mp.BaseMessage[T], to *T) {
	atomic.StoreInt32(&l.messagesSent, atomic.LoadInt32(&l.messagesSent)+1)

	l.mtMu.Lock()
	defer l.mtMu.Unlock()

	l.messageTracks[m.Id] = append(l.messageTracks[m.Id], _Record[T]{from: m.From, to: to, reSend: m.ReSends})
}

func (l *Logger[T]) String(phoneBook map[*T]int) string {
	str := fmt.Sprint(
		fmt.Sprintf("Messages sent : %v\n", atomic.LoadInt32(&l.messagesSent)),
		"Displaying tracks :\n",
	)

	l.mtMu.RLock()
	defer l.mtMu.RUnlock()

	str += fmt.Sprintf("%v\n", l.messageTracks[0])

	for key, arr := range l.messageTracks {
		str += fmt.Sprintf("Track #%d, total messages - %d :\n", key, len(arr))

		maxResends := 0
		for _, message := range arr {
			maxResends = max(maxResends, message.reSend)
		}

		sortResends := make([][]int, maxResends+2)

		for _, message := range arr {
			sortResends[message.reSend+1] = append(sortResends[message.reSend+1], phoneBook[message.to])
		}

		for i, messageArr := range sortResends {
			for _, nodeNumber := range messageArr {
				str += fmt.Sprintf(" %d", nodeNumber)
			}

			if i+1 == len(sortResends) {
				str += "\n"
			} else {
				str += " ->"
			}
		}
	}

	return str
}

func NewLogger[T any]() *Logger[T] {
	return &Logger[T]{messageTracks: make(map[int][]_Record[T])}
}
