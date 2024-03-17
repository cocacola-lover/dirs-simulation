package fundamentals

type IMessage interface {
	Id() int
	From() INode
	Key() string
	IsValid() bool
	Resends() int
	Resend(from INode) IMessage
}
