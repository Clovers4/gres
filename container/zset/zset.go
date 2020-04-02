package zset

import (
	"bytes"
	"encoding/binary"
	"github.com/clovers4/gres/container/zset/skiplist"
	"github.com/clovers4/gres/util"
)

var DefaultByteOrder = binary.LittleEndian

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

func (zs *ZSet) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// write total. the total must > 0
	total := zs.Length()
	if err := binary.Write(buf, DefaultByteOrder, int64(total)); err != nil {
		return nil, err
	}

	// loop write score and val
	for n := zs.Rank(1); n != nil; n = n.Next() {
		if err := binary.Write(buf, DefaultByteOrder, n.Score()); err != nil {
			return nil, err
		}

		val := util.StringToBytes(n.Val())
		if err := binary.Write(buf, DefaultByteOrder, len(val)); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, DefaultByteOrder, val); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (zs *ZSet) Unmarshal(b []byte) {

}
