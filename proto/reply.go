package proto

import "fmt"

type ReplyKind uint8

const (
	ReplyKindStatus = iota
	ReplyKindErr
	ReplyKindInt
	ReplyKindBlukString
	ReplyKindArrays
)

type Reply struct {
	Kind ReplyKind
	Val  interface{}
	Err  error
}

func NewReply(kind ReplyKind, val interface{}, err error) *Reply {
	if err != nil {
		kind = ReplyKindErr
	}
	return &Reply{
		Kind: kind,
		Val:  val,
		Err:  err,
	}
}

// only for test
func (r *Reply) String() string {
	var kind string
	switch r.Kind {
	case ReplyKindStatus:
		kind = "Status"
	case ReplyKindErr:
		kind = "Err"
	case ReplyKindInt:
		kind = "Int"
	case ReplyKindBlukString:
		kind = "BlukString"
	case ReplyKindArrays:
		kind = "Arrays"
	}
	return fmt.Sprintf("[%v] val=%v, err=%v", kind, r.Val, r.Err)
}
