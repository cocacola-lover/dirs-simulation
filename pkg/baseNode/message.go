package basenode

type IMessage interface {
	Id() int
	From() *BaseNode
	Key() string
	IsValid() bool
	Resend(from *BaseNode) any
}

func Resend(m IMessage, from *BaseNode) IMessage {
	return m.Resend(from).(IMessage)
}
