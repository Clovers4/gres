package zset

import (
	"encoding/binary"
	"fmt"
	"github.com/clovers4/gres/engine/object/zset/skiplist"
	"io"

	"github.com/clovers4/gres/util"
)

var DefaultByteOrder = binary.BigEndian

// effective, so dont support concurrent ops.
type ZSet struct {
	m        map[string]float64
	skiplist *skiplist.Skiplist
}

func New() *ZSet {
	return &ZSet{
		m:        make(map[string]float64),
		skiplist: skiplist.New(),
	}
}

func (zs *ZSet) Add(score float64, member string) bool {
	// found
	if curScore, ok := zs.m[member]; ok {
		if curScore != score {
			zs.m[member] = score
			zs.skiplist.UpdateScore(curScore, score, member)
		}
		return false
	}
	// not found
	zs.m[member] = score
	zs.skiplist.Insert(score, member)
	return true
}

func (zs *ZSet) Incr(member string) {
	// found
	if curScore, ok := zs.m[member]; ok {
		zs.m[member]++
		zs.skiplist.UpdateScore(curScore, curScore+1, member)
	}
	// not found
	zs.m[member] = 1
	zs.skiplist.Insert(1, member)
}

func (zs *ZSet) Delete(member string) (float64, bool) {
	// found
	if curScore, ok := zs.m[member]; ok {
		delete(zs.m, member)
		zs.skiplist.Delete(curScore, member)
		return curScore, true
	}
	return 0, false
}

func (zs *ZSet) Get(member string) (float64, bool) {
	score, ok := zs.m[member]
	return score, ok
}

func (zs *ZSet) GetRankByMember(member string) (rank int, existed bool) {
	score, existed := zs.m[member]
	if !existed {
		return -1, false
	}
	return zs.skiplist.GetRankByScore(score, &member)
}

func (zs *ZSet) GetNodeByRank(rank int) *skiplist.SkiplistNode {
	return zs.skiplist.GetNodeByRank(rank)
}

func (zs *ZSet) Length() int {
	return zs.skiplist.Length() // can also use len(zs.m), but maybe skiplist.Length() is more fast
}

// Only for test
func (zs *ZSet) String() string {
	var s string
	s += "{"
	for n := zs.skiplist.Front(); n != nil; n = n.Next() {
		s += fmt.Sprintf("%v : %v, ", n.Val(), n.Score())
	}

	if len(s) > 2 {
		s = s[:len(s)-2]
	}
	s += "}"
	return s
}

func (zs *ZSet) Marshal(w io.Writer) error {
	// write total. the total must > 0
	total := zs.Length()
	if err := util.Write(w, int64(total)); err != nil {
		return err
	}

	// loop write score and val
	for n := zs.GetNodeByRank(0); n != nil; n = n.Next() {
		if err := util.Write(w, n.Score()); err != nil {
			return err
		}

		if err := util.Write(w, n.Val()); err != nil {
			return err
		}
	}
	return nil
}

func (zs *ZSet) Unmarshal(r io.Reader) error {
	var total int64
	if err := util.Read(r, &total); err != nil {
		return err
	}

	for i := 0; i < int(total); i++ {
		var score float64
		if err := util.Read(r, &score); err != nil {
			return err
		}

		var val string
		if err := util.Read(r, &val); err != nil {
			return err
		}
		zs.Add(score, val)
	}
	return nil
}
