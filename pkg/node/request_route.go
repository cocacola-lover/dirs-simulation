package node

type _RouteRequest struct {
	id       int
	key      string
	from     INode
	sentTo   []INode
	routedTo INode
}

func _NewRouteRequest(id int, key string, from INode) _RouteRequest {
	return _RouteRequest{id: id, key: key, from: from}
}
