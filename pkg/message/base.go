package message

import (
	fp "dirs/simulation/pkg/fundamentals"
	"sync/atomic"
)

var id int32 = 0

func getId() int32 {
	lid := atomic.LoadInt32(&id)
	defer atomic.StoreInt32(&id, lid+1)
	return lid
}

type BaseMessage struct {
	id      int32
	key     string
	reSends int
	from    fp.INode
}

func (m BaseMessage) Id() int {
	return int(m.id)
}
func (m BaseMessage) Key() string {
	return m.key
}
func (m BaseMessage) From() fp.INode {
	return m.from
}
func (m BaseMessage) Resends() int {
	return m.reSends
}

func (m BaseMessage) Done() {}

func (m BaseMessage) Resend(from fp.INode) fp.IMessage {
	m.reSends++
	m.from = from

	return m
}

func NewBaseMessage(key string, from fp.INode) BaseMessage {
	return BaseMessage{id: getId(), key: key, reSends: -1, from: from}
}
