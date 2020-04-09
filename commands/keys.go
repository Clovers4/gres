package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/errs"
	"github.com/clovers4/gres/proto"
	"github.com/clovers4/gres/util"
)

// KEYS
func init() {
	registerCmd("quit", 1, quitCmd)
	registerCmd("expire", 3, expireCmd)
	registerCmd("ttl", 2, ttlCmd)
	registerCmd("dbsize", 1, dbsizeCmd)
	registerCmd("exists", 2, existsCmd)
	registerCmd("del", 2, delCmd)
	registerCmd("type", 2, typeCmd)
	registerCmd("keys", 2, keysCmd)
}

func quitCmd(db *engine.DB, args []string) *proto.Reply {
	return proto.NewReply(proto.ReplyKindStatus, "OK", nil)
}

func expireCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	secondsS := args[2]
	var seconds int
	var err error

	if seconds, err = util.String2Int(secondsS); err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	if db.Expire(key, seconds) {
		return proto.NewReply(proto.ReplyKindInt, 1, nil)
	}
	return proto.NewReply(proto.ReplyKindInt, 0, nil)
}

func ttlCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]

	ttl := db.Ttl(key)
	return proto.NewReply(proto.ReplyKindInt, ttl, nil)
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

func keysCmd(db *engine.DB, args []string) *proto.Reply {
	pattern := args[1]
	ks, err := db.Keys(pattern)
	is := make([]interface{}, len(ks))
	for i := 0; i < len(ks); i++ {
		is[i] = ks[i]
	}
	return proto.NewReply(proto.ReplyKindArrays, is, err)
}
