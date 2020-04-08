package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/proto"
	"github.com/clovers4/gres/util"
)

// SET
func init() {
	registerCmd("sadd", -3, saddCmd)
	registerCmd("srem", -3, sremCmd)
	registerCmd("scard", 2, scardCmd)
	registerCmd("sismember", 3, sismemberCmd)
	registerCmd("smembers", 2, smembersCmd)
	registerCmd("sinter", -3, sinterCmd)
	registerCmd("sunion", -3, sunionCmd)
	registerCmd("sdiff", -3, sdiffCmd)
}

func saddCmd(db *engine.DB, args []string) *proto.Reply {
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

	count, err := db.SAdd(key, vs...)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func sremCmd(db *engine.DB, args []string) *proto.Reply {
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

	count, err := db.SRem(key, vs...)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func scardCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	count, err := db.SCard(key)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func sismemberCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	val := args[2]
	var found bool
	var err error
	if valNum, ok := util.String2Num(val); ok {
		found, err = db.SIsMember(key, valNum)
	}

	found, err = db.SIsMember(key, val)
	if found {
		return proto.NewReply(proto.ReplyKindInt, 1, err)
	}
	return proto.NewReply(proto.ReplyKindInt, 0, err)
}

func smembersCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	vals, err := db.SMembers(key)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}

func sinterCmd(db *engine.DB, args []string) *proto.Reply {
	vals, err := db.SInter(args[1:]...)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}

func sunionCmd(db *engine.DB, args []string) *proto.Reply {
	vals, err := db.SUnion(args[1:]...)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}

func sdiffCmd(db *engine.DB, args []string) *proto.Reply {
	vals, err := db.SDiff(args[1:]...)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}
