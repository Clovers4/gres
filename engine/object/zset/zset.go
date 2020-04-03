package zset

import (
	"encoding/binary"
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

func (zs *ZSet) Add(score float64, val string) {
	// found
	if curScore, ok := zs.m[val]; ok {
		if curScore != score {
			zs.m[val] = score
			zs.skiplist.UpdateScore(curScore, score, val)
		}
		return
	}
	// not found
	zs.m[val] = score
	zs.skiplist.Insert(score, val)
}

func (zs *ZSet) Incr(val string) {
	// found
	if curScore, ok := zs.m[val]; ok {
		zs.m[val]++
		zs.skiplist.UpdateScore(curScore, curScore+1, val)
	}
	// not found
	zs.m[val] = 1
	zs.skiplist.Insert(1, val)
}

func (zs *ZSet) Delete(val string) {
	// found
	if curScore, ok := zs.m[val]; ok {
		delete(zs.m, val)
		zs.skiplist.Delete(curScore, val)
	}
}

func (zs *ZSet) Get(val string) (float64, bool) {
	score, ok := zs.m[val]
	return score, ok
}

func (zs *ZSet) Rank(rank int) *skiplist.SkiplistNode {
	return zs.skiplist.GetNodeByRank(rank)
}

func (zs *ZSet) Length() int {
	return zs.skiplist.Length() // can also use len(zs.m), but maybe skiplist.Length() is more fast
}

func (zs *ZSet) Marshal(w io.Writer) error {
	// write total. the total must > 0
	total := zs.Length()
	if err := util.Write(w, int64(total)); err != nil {
		return err
	}

	// loop write score and val
	for n := zs.Rank(1); n != nil; n = n.Next() {
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
