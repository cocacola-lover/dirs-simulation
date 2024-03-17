package logger

import (
	fp "dirs/simulation/pkg/fundamentals"
	"fmt"
	"sync"
	"sync/atomic"
)

// # The purpose of the logger is to keep track of operations for statistic's sake.
//
// ## The logger should count the number of messages sent between nodes.
// ## The logger should keep track of searches in the network and their results.
type Logger struct {
	messagesSent int32

	messageTracks map[int][]_Record[fp.INode]
	mtMu          sync.RWMutex
}

func (l *Logger) MessagesSent() int {
	return int(atomic.LoadInt32(&l.messagesSent))
}

func (l *Logger) AddMessage(m fp.IMessage, to fp.INode) {
	atomic.StoreInt32(&l.messagesSent, atomic.LoadInt32(&l.messagesSent)+1)

	l.mtMu.Lock()
	defer l.mtMu.Unlock()

	l.messageTracks[m.Id()] = append(l.messageTracks[m.Id()], _Record[fp.INode]{from: m.From(), to: to, reSend: m.Resends()})
}

func (l *Logger) String(phoneBook map[fp.INode]int) string {
	str := fmt.Sprint(
		fmt.Sprintf("Messages sent : %v\n", atomic.LoadInt32(&l.messagesSent)),
		"Displaying tracks :\n",
	)

	l.mtMu.RLock()
	defer l.mtMu.RUnlock()

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

func NewLogger() *Logger {
	return &Logger{messageTracks: make(map[int][]_Record[fp.INode])}
}
