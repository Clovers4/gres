package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/proto"
)

var (
	ErrUnknownCmd     = errors.New("ERR unknown command")
	ErrWrongNumArgs   = errors.New("ERR wrong number of arguments for the command")
	ErrWrongTypeInt   = errors.New("ERR value is not an integer")
	ErrInvalidDbIndex = errors.New("ERR invalid DB index")
)

// commands should be read-only
var commands map[string]Command

func registerCmd(name string, arity int, do doFunc) {
	if commands == nil {
		commands = make(map[string]Command)
	}
	if commands[name] != nil {
		panic(fmt.Errorf("cmd %s is already registerd", name))
	}
	commands[name] = &cmd{
		name,
		arity,
		do,
	}
}

func GetCmd(name string) Command {
	return commands[name]
}

type Command interface {
	Do(ctx context.Context, args []string) *proto.Reply
}

type doFunc func(db *engine.DB, args []string) *proto.Reply

type cmd struct {
	name  string
	arity int // Number of arguments, it is possible to use -N to say >= N
	do    doFunc
}

func (c *cmd) Do(ctx context.Context, args []string) *proto.Reply {
	if c.arity > 0 && len(args) != c.arity ||
		len(args) < -c.arity {
		return proto.NewReply(proto.ReplyKindErr, nil, ErrWrongNumArgs)
	}
	db := engine.CtxGetDB(ctx)
	return c.do(db, args)
}
