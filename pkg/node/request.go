package node

type Request struct {
	id           int
	key          string
	from         INode
	sentTo       []INode
	routedTo     []INode
	awaitingFrom INode
}

func (r Request) Id() int {
	return r.id
}

func _NewRequest(id int, key string, from INode) Request {
	return Request{id: id, key: key, from: from}
}
