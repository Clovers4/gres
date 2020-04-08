package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/errs"
	"github.com/clovers4/gres/proto"
	"github.com/clovers4/gres/util"
)

// PLAIN
func init() {
	registerCmd("set", 3, setCmd) // todo:
	registerCmd("get", 2, getCmd)
	registerCmd("getset", 3, getsetCmd)

	registerCmd("incr", 2, incrCmd)
	registerCmd("incrby", 3, incrbyCmd)
	registerCmd("decr", 2, decrCmd)
	registerCmd("decrby", 3, decrbyCmd)
}

func setCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val := args[2]

	if valNum, ok := util.String2Num(val); ok {
		err := db.Set(key, valNum)
		return proto.NewReply(proto.ReplyKindStatus, "OK", err)
	}
	err := db.Set(key, val)
	return proto.NewReply(proto.ReplyKindStatus, "OK", err)
}

func getCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]

	val, err := db.Get(key)
	return proto.NewReply(proto.ReplyKindBlukString, val, err)
}

func getsetCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val := args[2]

	if valNum, ok := util.String2Num(val); ok {
		oldVal, err := db.GetSet(key, valNum)
		return proto.NewReply(proto.ReplyKindBlukString, oldVal, err)
	}
	oldVal, err := db.GetSet(key, val)
	return proto.NewReply(proto.ReplyKindBlukString, oldVal, err)
}

func incrCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val, err := db.Incr(key)
	return proto.NewReply(proto.ReplyKindInt, val, err)
}

func incrbyCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val := args[2]

	valInt, err := util.String2Int(val)
	if err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	afterVal, err := db.IncrBy(key, valInt)
	return proto.NewReply(proto.ReplyKindInt, afterVal, err)
}

func decrCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val, err := db.Decr(key)
	return proto.NewReply(proto.ReplyKindInt, val, err)
}

func decrbyCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val := args[2]

	valInt, err := util.String2Int(val)
	if err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	afterVal, err := db.DecrBy(key, valInt)
	return proto.NewReply(proto.ReplyKindInt, afterVal, err)
}
