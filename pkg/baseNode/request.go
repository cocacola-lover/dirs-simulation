package basenode

import fp "dirs/simulation/pkg/fundamentals"

type _Request struct {
	fp.IMessage
	SentTo []fp.INode
}

func (r _Request) Done(by fp.INode) {
	r.stopSearch(by)
	r.IMessage.Done(by)
}

func (r _Request) stopSearch(by fp.INode) {
	for _, n := range r.SentTo {
		n.StopSearch(r.Id(), by)
	}
}

func _NewRequest(m fp.IMessage, sendTo []fp.INode) _Request {
	return _Request{IMessage: m, SentTo: sendTo}
}
