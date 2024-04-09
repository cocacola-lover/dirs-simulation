package loggernode

import (
	lp "dirs/simulation/pkg/nLogger"
	"dirs/simulation/pkg/node"
	snp "dirs/simulation/pkg/searcherNode"
	"sync"
)

type LoggerNode struct {
	*snp.SearcherNode
	logger *lp.Logger

	searches     map[int]chan bool
	searchesLock sync.Mutex
}

func (ln *LoggerNode) ReceiveRouteMessage(id int, key string, from node.INode) bool {
	hasAccepted := ln.INode.ReceiveRouteMessage(id, key, from)

	if hasAccepted {
		go ln.logger.AddRouteMessageReceive(id, from, ln)
	} else {
		go ln.logger.AddDeniedRouteMessageReceive(id, from, ln)
	}

	return hasAccepted
}

func (ln *LoggerNode) ConfirmRouteMessage(id int, from node.INode) bool {
	hasAccepted := ln.INode.ConfirmRouteMessage(id, from)

	if hasAccepted {
		go ln.logger.AddRouteMessageConfirm(id, from, ln)
	} else {
		go ln.logger.AddDeniedRouteMessageConfirm(id, from, ln)
	}

	return hasAccepted
}

func (ln *LoggerNode) ReceiveDownloadMessage(id int, key string, from node.INode) {
	ln.INode.ReceiveDownloadMessage(id, key, from)
}

func (ln *LoggerNode) ConfirmDownloadMessage(id int, val string, from node.INode) {
	ln.INode.ConfirmDownloadMessage(id, val, from)
	go ln.logger.AddDownloadMessage(id, from, ln.GetSelfAddress())
}

func (ln *LoggerNode) StartSearchAndWatch(key string) int {

	ch := make(chan bool)
	id := ln.SearcherNode.StartSearchAndWatch(key, ch)

	ln.searchesLock.Lock()
	ln.searches[id] = ch
	ln.searchesLock.Unlock()

	ln.logger.StartSearch(id, ln.GetSelfAddress())

	return id
}

func (ln *LoggerNode) PutVal(key, val string) {
	ln.SearcherNode.PutVal(key, val)

	ln.searchesLock.Lock()
	defer ln.searchesLock.Unlock()
	for id, ch := range ln.searches {
		select {
		case <-ch:
			delete(ln.searches, id)
			ln.logger.EndSearch(id)
		default:
			continue
		}
	}
}

func (ln *LoggerNode) WaitToFinishAllSearches(waitGroup *sync.WaitGroup) {
	for {
		ln.searchesLock.Lock()
		if len(ln.searches) == 0 {
			ln.searchesLock.Unlock()
			break
		}

		var pickCh chan bool
		for _, ch := range ln.searches {
			pickCh = ch
			break
		}
		ln.searchesLock.Unlock()

		<-pickCh
	}
	if waitGroup != nil {
		waitGroup.Done()
	}
}

func (ln *LoggerNode) RetryMessages(rs []node.Request) []int {
	newIds := ln.SearcherNode.RetryMessages(rs)

	oldIds := make([]int, 0, len(rs))
	for _, rs := range rs {
		oldIds = append(oldIds, rs.Id())
	}

	ln.logger.ChangeIdForMessages(oldIds, newIds)

	return newIds
}

func (ln *LoggerNode) Fail() {
	ln.SearcherNode.Fail()
	go ln.logger.AddFailedNode()
}

func (ln *LoggerNode) ReceiveFaultMessage(from node.INode, about []int) []int {
	go ln.logger.AddFaultMessageReceive()
	return ln.SearcherNode.ReceiveFaultMessage(from, about)
}

func (ln *LoggerNode) Close() {
	ln.SearcherNode.Close()

	ln.searchesLock.Lock()
	defer ln.searchesLock.Unlock()

	ln.searches = nil
}

func NewLoggerNode(node *snp.SearcherNode, logger *lp.Logger) *LoggerNode {
	ln := &LoggerNode{SearcherNode: node, logger: logger, searches: make(map[int]chan bool)}
	ln.SetSelfAddress(ln)
	return ln
}
