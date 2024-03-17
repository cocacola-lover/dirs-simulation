package message

import (
	fp "dirs/simulation/pkg/fundamentals"
)

type FirstMessage struct {
	BaseMessage

	closeCh chan bool
}

func (m FirstMessage) Done() {
	close(m.closeCh)
}

func (m FirstMessage) WaitForDone() {
	<-m.closeCh
}

func NewFirstMessage(key string, from fp.INode) FirstMessage {
	return FirstMessage{closeCh: make(chan bool), BaseMessage: BaseMessage{id: getId(), key: key, reSends: -1, from: from}}
}
