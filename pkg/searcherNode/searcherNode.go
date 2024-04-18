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

func (sn *SearcherNode) PutKey(key, val string) {
	sn.INode.PutKey(key, val)

	sn.searchesLock.Lock()
	defer sn.searchesLock.Unlock()

	v, ok := sn.searches[key]

	if ok {
		delete(sn.searches, key)
		close(v)
	}
}

func (sn *SearcherNode) Close() {
	sn.INode.Close()

	sn.searchesLock.Lock()
	defer sn.searchesLock.Unlock()

	for key, ch := range sn.searches {
		close(ch)
		delete(sn.searches, key)
	}
}

func NewSearchNode(bn node.INode) *SearcherNode {
	n := &SearcherNode{INode: bn, searches: make(map[string]chan bool)}
	n.SetSelfAddress(n)
	return n
}
