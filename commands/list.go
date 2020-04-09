package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/errs"
	"github.com/clovers4/gres/proto"
	"github.com/clovers4/gres/util"
)

// LIST
func init() {
	registerCmd("lpush", -2, lpushCmd)
	registerCmd("rpush", -3, rpushCmd)
	registerCmd("lpop", 2, lpopCmd)
	registerCmd("rpop", 2, rpopCmd)

	registerCmd("lrange", 4, lrangeCmd)
	registerCmd("lindex", 3, lindexCmd)
	registerCmd("lset", 4, lsetCmd)
}

func lpushCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	vals := args[2:]

	var vs []interface{}
	for _, val := range vals {
		if valNum, ok := util.String2Num(val); ok {
			vs = append(vs, valNum)
		} else {
			vs = append(vs, val)
		}
	}

	len, err := db.LPush(key, vs...)
	return proto.NewReply(proto.ReplyKindInt, len, err)
}

func rpushCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	vals := args[2:]

	var vs []interface{}
	for _, val := range vals {
		if valNum, ok := util.String2Num(val); ok {
			vs = append(vs, valNum)
		} else {
			vs = append(vs, val)
		}
	}

	len, err := db.RPush(key, vs...)
	return proto.NewReply(proto.ReplyKindInt, len, err)
}

func lpopCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]

	oldVal, err := db.LPop(key)
	return proto.NewReply(proto.ReplyKindBlukString, oldVal, err)
}

func rpopCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]

	oldVal, err := db.RPop(key)
	return proto.NewReply(proto.ReplyKindBlukString, oldVal, err)
}

func lrangeCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	startS := args[2]
	endS := args[3]

	start, err := util.String2Int(startS)
	if err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	end, err := util.String2Int(endS)
	if err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	vals, err := db.LRange(key, start, end)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}

func lindexCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	indexS := args[2]

	index, err := util.String2Int(indexS)
	if err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	val, err := db.LIndex(key, index)
	return proto.NewReply(proto.ReplyKindBlukString, val, err)
}

func lsetCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	indexS := args[2]
	valS := args[3]

	index, err := util.String2Int(indexS)
	if err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	valNum, ok := util.String2Num(valS)
	if ok {
		_, err := db.LSet(key, index, valNum)
		return proto.NewReply(proto.ReplyKindStatus, "OK", err)
	}

	_, err = db.LSet(key, index, valS)
	return proto.NewReply(proto.ReplyKindStatus, "OK", err)
}
