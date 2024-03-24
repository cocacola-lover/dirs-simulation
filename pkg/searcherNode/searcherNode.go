package searchernode

import (
	"dirs/simulation/pkg/node"
	"sync"
)

type SearcherNode struct {
	node.INode

	searches     map[string]chan bool
	searchesLock sync.RWMutex
}

func (sn *SearcherNode) StartSearch(id int, key string, ch chan bool) {
	sn.searchesLock.Lock()
	sn.searches[key] = ch
	sn.searchesLock.Unlock()

	if !sn.INode.ReceiveRouteMessage(id, key, sn.GetSelfAddress()) {
		panic("ReceiveRouteMessage on start search returned false")
	}
}

func (sn *SearcherNode) PutVal(key, val string) {
	sn.INode.PutVal(key, val)

	sn.searchesLock.RLock()
	close(sn.searches[key])
	sn.searchesLock.RUnlock()
}

func NewSearchNode(bn node.INode) *SearcherNode {
	n := &SearcherNode{INode: bn, searches: make(map[string]chan bool)}
	n.SetSelfAddress(n)
	return n
}
