package node

type _RouteRequest struct {
	id       int
	key      string
	from     *Node
	sentTo   []*Node
	routedTo *Node
}

func _NewRouteRequest(id int, key string, from *Node) _RouteRequest {
	return _RouteRequest{id: id, key: key, from: from}
}
