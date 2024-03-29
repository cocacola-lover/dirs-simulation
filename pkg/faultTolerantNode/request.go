package faulttolerantnode

type _Request struct {
	id           int
	key          string
	from         *Node
	sentTo       []*Node
	routedTo     []*Node
	awaitingFrom *Node
}

func _NewRequest(id int, key string, from *Node) _Request {
	return _Request{id: id, key: key, from: from}
}
