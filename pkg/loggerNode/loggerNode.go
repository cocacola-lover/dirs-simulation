package loggernode

import (
	lp "dirs/simulation/pkg/nLogger"
	"dirs/simulation/pkg/node"
)

type LoggerNode struct {
	node.INode
	logger *lp.Logger
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

func (ln *LoggerNode) TimeoutRouteMessage(id int, from node.INode) bool {
	hasAccepted := ln.INode.TimeoutRouteMessage(id, from)

	if hasAccepted {
		go ln.logger.AddRouteMessageTimeout(id, from, ln)
	} else {
		go ln.logger.AddDeniedRouteMessageTimeout(id, from, ln)
	}

	return hasAccepted
}

func (ln *LoggerNode) ReceiveDownloadMessage(id int, key string, from node.INode) {
	ln.INode.ReceiveDownloadMessage(id, key, from)
}

func (ln *LoggerNode) ConfirmDownloadMessage(id int, val string, from node.INode) {
	ln.INode.ConfirmDownloadMessage(id, val, from)
	go ln.logger.AddDownloadMessage(id, from)
}

func NewLoggerNode(node node.INode, logger *lp.Logger) *LoggerNode {
	ln := &LoggerNode{INode: node, logger: logger}
	ln.INode.SetSelfAddress(ln)
	return ln
}
