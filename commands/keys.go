package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/proto"
)

// KEYS
func init() {
	registerCmd("dbsize", 1, dbsizeCmd)
	registerCmd("exists", 2, existsCmd)
	registerCmd("del", 2, delCmd)
	registerCmd("type", 2, typeCmd)
}

func dbsizeCmd(db *engine.DB, args []string) *proto.Reply {
	len := db.DbSize()
	return proto.NewReply(proto.ReplyKindInt, len, nil)
}

func existsCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	ok := db.Exists(key)
	if ok {
		return proto.NewReply(proto.ReplyKindInt, 1, nil)
	}
	return proto.NewReply(proto.ReplyKindInt, 0, nil)
}

func delCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1:]
	count := db.Del(key...)
	return proto.NewReply(proto.ReplyKindInt, count, nil)
}

func typeCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	t := db.Type(key)
	return proto.NewReply(proto.ReplyKindBlukString, t, nil)
}
