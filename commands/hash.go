package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/errs"
	"github.com/clovers4/gres/proto"
	"github.com/clovers4/gres/util"
)

// ZSET
func init() {
	registerCmd("hset", 4, hsetCmd) //todo
	registerCmd("hget", 3, hgetCmd)
	registerCmd("hdel", -3, hdelCmd)
	registerCmd("hlen", 2, hlenCmd)
	registerCmd("hexists", 3, hexistsCmd)
	registerCmd("hkeys", 2, hkeysCmd)
	registerCmd("hvals", 2, hvalsCmd)
	registerCmd("hgetall", 2, hgetallCmd)
	registerCmd("hincrby", 4, hincrbyCmd)
}

func hsetCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	field := args[2]
	valS := args[3]

	if val, ok := util.String2Num(valS); ok {
		count, err := db.HSet(key, field, val)
		return proto.NewReply(proto.ReplyKindInt, count, err)
	}

	count, err := db.HSet(key, field, valS)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func hgetCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	field := args[2]

	val, err := db.HGet(key, field)
	return proto.NewReply(proto.ReplyKindBlukString, val, err)
}

func hdelCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	fields := args[2:]

	count, err := db.HDel(key, fields...)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func hlenCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	count, err := db.HLen(key)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func hexistsCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	field := args[2]
	found, err := db.HExists(key, field)
	if found {
		return proto.NewReply(proto.ReplyKindInt, 1, err)
	}
	return proto.NewReply(proto.ReplyKindInt, 0, err)
}

func hkeysCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	vals, err := db.HKeys(key)

	is := make([]interface{}, len(vals))
	for i := range vals {
		is[i] = vals[i]
	}
	return proto.NewReply(proto.ReplyKindArrays, is, err)
}

func hvalsCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	vals, err := db.HVals(key)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}

func hgetallCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	vals, err := db.HGetAll(key)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}

func hincrbyCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	field := args[2]
	incrementS := args[3]
	var increment int
	var err error
	if increment, err = util.String2Int(incrementS); err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	val, err := db.HIncrBy(key, field, increment)
	return proto.NewReply(proto.ReplyKindInt, val, err)
}
