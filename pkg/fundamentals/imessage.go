package fundamentals

type IMessage interface {
	Id() int
	From() INode
	Key() string
	Done(by INode)
	Resends() int
	Resend(from INode) IMessage
}
