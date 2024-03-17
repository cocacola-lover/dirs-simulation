package logger

type _Record[T any] struct {
	from   T
	to     T
	reSend int
}
