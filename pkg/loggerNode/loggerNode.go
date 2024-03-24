package loggernode

import (
	lp "dirs/simulation/pkg/nLogger"
	"dirs/simulation/pkg/node"
)

type ComposobleNode interface {
	node.INode
	SetSelfAddress(n node.INode)
}

type LoggerNode struct {
	ComposobleNode
	logger *lp.Logger
}

func (ln *LoggerNode) ReceiveRouteMessage(id int, key string, from node.INode) bool {
	hasAccepted := ln.ComposobleNode.ReceiveRouteMessage(id, key, from)

	if hasAccepted {
		go ln.logger.AddRouteMessageReceive(id, from, ln)
	} else {
		go ln.logger.AddDeniedRouteMessageReceive(id, from, ln)
	}

	return hasAccepted
}

func (ln *LoggerNode) ConfirmRouteMessage(id int, from node.INode) bool {
	hasAccepted := ln.ComposobleNode.ConfirmRouteMessage(id, from)

	if hasAccepted {
		go ln.logger.AddRouteMessageConfirm(id, from, ln)
	} else {
		go ln.logger.AddDeniedRouteMessageConfirm(id, from, ln)
	}

	return hasAccepted
}

func (ln *LoggerNode) TimeoutRouteMessage(id int, from node.INode) bool {
	hasAccepted := ln.ComposobleNode.TimeoutRouteMessage(id, from)

	if hasAccepted {
		go ln.logger.AddRouteMessageTimeout(id, from, ln)
	} else {
		go ln.logger.AddDeniedRouteMessageTimeout(id, from, ln)
	}

	return hasAccepted
}

func (ln *LoggerNode) ReceiveDownloadMessage(id int, key string, from node.INode) {
	ln.ComposobleNode.ReceiveDownloadMessage(id, key, from)
}

func (ln *LoggerNode) ConfirmDownloadMessage(id int, val string, from node.INode) {
	ln.ComposobleNode.ConfirmDownloadMessage(id, val, from)
	go ln.logger.AddDownloadMessage(id, from)
}

func NewLoggerNode(node ComposobleNode, logger *lp.Logger) *LoggerNode {
	ln := &LoggerNode{ComposobleNode: node, logger: logger}
	ln.ComposobleNode.SetSelfAddress(ln)
	return ln
}
