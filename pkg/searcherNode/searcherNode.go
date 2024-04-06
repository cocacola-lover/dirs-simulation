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

func (sn *SearcherNode) StartSearchAndWatch(key string, ch chan bool) int {
	sn.searchesLock.Lock()
	sn.searches[key] = ch
	sn.searchesLock.Unlock()

	return sn.GetSelfAddress().StartSearch(key)
}

func (sn *SearcherNode) PutVal(key, val string) {
	sn.INode.PutVal(key, val)

	sn.searchesLock.RLock()
	v, ok := sn.searches[key]
	if !ok {
		panic("Don't have a channel to close in PutVal")
	}
	close(v)
	sn.searchesLock.RUnlock()
}

func NewSearchNode(bn node.INode) *SearcherNode {
	n := &SearcherNode{INode: bn, searches: make(map[string]chan bool)}
	n.SetSelfAddress(n)
	return n
}
