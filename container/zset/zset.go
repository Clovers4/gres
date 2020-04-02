package zset

import (
	"encoding/binary"
	"fmt"
	"github.com/clovers4/gres/container/zset/skiplist"
	"github.com/clovers4/gres/util"
	"io"
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

func (zs *ZSet) Rank(rank int) *skiplist.SkiplistNode {
	return zs.skiplist.GetNodeByRank(rank)
}

func (zs *ZSet) Length() int {
	return zs.skiplist.Length() // can also use len(zs.m), but maybe skiplist.Length() is more fast
}

func (zs *ZSet) Marshal(w io.Writer) error {
	// write total. the total must > 0
	total := zs.Length()
	if err := binary.Write(w, DefaultByteOrder, int64(total)); err != nil {
		return err
	}

	// loop write score and val
	for n := zs.Rank(1); n != nil; n = n.Next() {
		if err := binary.Write(w, DefaultByteOrder, n.Score()); err != nil {
			return err
		}

		val := util.StringToBytes(n.Val())
		if err := binary.Write(w, DefaultByteOrder, int64(len(val))); err != nil {
			return err
		}
		if err := binary.Write(w, DefaultByteOrder, val); err != nil {
			return err
		}
	}
	return nil
}

func (zs *ZSet) Unmarshal(r io.Reader) error {
	var total int64
	if err := binary.Read(r, DefaultByteOrder, &total); err != nil {
		return err
	}

	fmt.Println("total", total)
	for i := 0; i < int(total); i++ {
		var score float64
		if err := binary.Read(r, DefaultByteOrder, &score); err != nil {
			return err
		}
		fmt.Println("score", score)

		var len int64
		if err := binary.Read(r, DefaultByteOrder, &len); err != nil {
			return err
		}

		bs := make([]byte, len)
		if _, err := io.ReadFull(r, bs); err != nil {
			return err
		}

		val := util.BytesToString(bs)
		zs.Add(score, val)
	}
	return nil
}
