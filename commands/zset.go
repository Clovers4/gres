package commands

import (
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/errs"
	"github.com/clovers4/gres/proto"
	"github.com/clovers4/gres/util"
	"strings"
)

// ZSET
func init() {
	registerCmd("zadd", 4, zaddCmd) //todo
	registerCmd("zcard", 2, zcardCmd)
	registerCmd("zscore", 3, zscoreCmd)
	registerCmd("zrank", 3, zrankCmd)
	registerCmd("zrem", -3, zremCmd)
	registerCmd("zincrby", 4, zincrbyCmd)
	registerCmd("zrange", -4, zrangeCmd)
}

func zaddCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	scoreS := args[2]
	member := args[3]

	var score float64
	var err error
	if score, err = util.String2Float(scoreS); err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotFloat)
	}

	count, err := db.ZAdd(key, score, member)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func zcardCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	count, err := db.ZCard(key)
	return proto.NewReply(proto.ReplyKindInt, count, err)
}

func zscoreCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	member := args[2]
	score, err := db.ZScore(key, member)
	if score == nil {
		return proto.NewReply(proto.ReplyKindBlukString, nil, err)
	}
	return proto.NewReply(proto.ReplyKindBlukString, score, err)
}

func zrankCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	member := args[2]
	rank, err := db.ZRank(key, member)

	if rank == nil {
		return proto.NewReply(proto.ReplyKindBlukString, nil, err)
	}
	return proto.NewReply(proto.ReplyKindInt, *rank, err)
}

func zremCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	members := args[2:]
	num, err := db.ZRem(key, members...)
	return proto.NewReply(proto.ReplyKindInt, num, err)
}

func zincrbyCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	scoreS := args[2]
	member := args[3]

	var score float64
	var err error
	if score, err = util.String2Float(scoreS); err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotFloat)
	}

	count, err := db.ZIncrBy(key, score, member)
	return proto.NewReply(proto.ReplyKindBlukString, count, err)
}

func zrangeCmd(db *engine.DB, args []string) *proto.Reply {
	key := args[1]
	startS := args[2]
	endS := args[3]
	var start int
	var end int
	var err error
	var withScore bool
	if len(args) == 5 && strings.ToLower(args[4]) == "withscores" {
		withScore = true
	}

	if start, err = util.String2Int(startS); err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}
	if end, err = util.String2Int(endS); err != nil {
		return proto.NewReply(proto.ReplyKindErr, nil, errs.ErrIsNotInt)
	}

	vals, err := db.ZRange(key, start, end, withScore)
	return proto.NewReply(proto.ReplyKindArrays, vals, err)
}
