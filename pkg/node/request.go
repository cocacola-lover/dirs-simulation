package node

type _Request struct {
	id           int
	key          string
	from         INode
	sentTo       []INode
	routedTo     []INode
	awaitingFrom INode
}

func _NewRequest(id int, key string, from INode) _Request {
	return _Request{id: id, key: key, from: from}
}
