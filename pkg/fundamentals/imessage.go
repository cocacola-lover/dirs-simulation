package fundamentals

type IMessage interface {
	Id() int
	From() INode
	Key() string
	Done()
	Resends() int
	Resend(from INode) IMessage
}
