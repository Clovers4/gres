package zset

import (
	"fmt"
	"io"
	"sync"

	"github.com/clovers4/gres/util"
	"github.com/clovers4/gres/zset/skiplist"
)

// this is used for expire

// effective, so dont support concurrent ops.
type ZSet struct {
	m        map[string]int64
	skiplist *skiplist.Skiplist

	sync.RWMutex
}

func New() *ZSet {
	return &ZSet{
		m:        make(map[string]int64),
		skiplist: skiplist.New(),
	}
}

func (zs *ZSet) AddZSet(z2 *ZSet) {
	for member, score := range z2.m {
		// ignore -1
		if score == -1 {
			zs.Delete(member)
		} else {
			zs.Add(score, member)
		}
	}
}

func (zs *ZSet) Add(score int64, member string) bool {
	zs.Lock()
	defer zs.Unlock()

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
	zs.Lock()
	defer zs.Unlock()

	// found
	if curScore, ok := zs.m[member]; ok {
		zs.m[member]++
		zs.skiplist.UpdateScore(curScore, curScore+1, member)
	}
	// not found
	zs.m[member] = 1
	zs.skiplist.Insert(1, member)
}

func (zs *ZSet) Delete(member string) (int64, bool) {
	zs.Lock()
	defer zs.Unlock()

	// found
	if curScore, ok := zs.m[member]; ok {
		delete(zs.m, member)
		zs.skiplist.Delete(curScore, member)
		return curScore, true
	}
	return 0, false
}

func (zs *ZSet) Get(member string) (int64, bool) {
	zs.RLock()
	defer zs.RUnlock()

	score, ok := zs.m[member]
	return score, ok
}

func (zs *ZSet) GetRankByMember(member string) (rank int, existed bool) {
	zs.RLock()
	defer zs.RUnlock()

	score, existed := zs.m[member]
	if !existed {
		return -1, false
	}
	return zs.skiplist.GetRankByScore(score)
}

// need rlock
func (zs *ZSet) GetNodeLeastScore(score int64) *skiplist.SkiplistNode {
	return zs.skiplist.GetNodeLeastScore(score)
}

// need rlock
func (zs *ZSet) GetNodeByRank(rank int) *skiplist.SkiplistNode {
	return zs.skiplist.GetNodeByRank(rank)
}

func (zs *ZSet) Length() int {
	zs.RLock()
	defer zs.RUnlock()
	return zs.skiplist.Length() // can also use len(zs.m), but maybe skiplist.Length() is more fast
}

// Only for test
func (zs *ZSet) String() string {
	zs.RLock()
	defer zs.RUnlock()

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
		var score int64
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
