package skiplist

import (
	"fmt"
	"testing"

	"github.com/clovers4/gres/engine/object/plain"
	"github.com/stretchr/testify/assert"
)

func array2String(vals []interface{}) string {
	var s string
	s += "{"

	for _, val := range vals {
		p := plain.New(val)
		s += p.String() + ", "
	}
	if len(s) > 2 {
		s = s[:len(s)-2]
	}
	s += "}"
	return s
}

func TestSkiplist(t *testing.T) {
	sl := New()
	sl.Insert(1, "a")
	sl.Insert(1, "a22")
	sl.Insert(3, "c")
	sl.Insert(2, "b")

	sl.Delete(1, "a22")

	var vals []interface{}
	for n := sl.Front(); n != nil; n = n.Next() {
		vals = append(vals, n.Score(), n.Val())
	}
	assert.Equal(t, 6, len(vals))
	assert.Equal(t, "{1, a, 2, b, 3, c}", array2String(vals))

	sl.UpdateScore(1, 5, "a")
	vals = []interface{}{}
	for n := sl.Front(); n != nil; n = n.Next() {
		vals = append(vals, n.Score(), n.Val())
	}
	assert.Equal(t, 6, len(vals))
	assert.Equal(t, "{2, b, 3, c, 5, a}", array2String(vals))

	target := sl.GetNodeByRank(0)
	assert.Equal(t, int64(2), target.Score())
	assert.Equal(t, "b", target.Val())

	target = sl.GetNodeByRank(1)
	assert.Equal(t, int64(3), target.Score())
	assert.Equal(t, "c", target.Val())

	target = sl.GetNodeByRank(5)
	assert.Nil(t, target)

	var rank int
	var existed bool
	rank, existed = sl.GetRankByScore(8)
	assert.Equal(t, false, existed)

	rank, existed = sl.GetRankByScore(5)
	assert.Equal(t, 2, rank)
	assert.Equal(t, true, existed)

	rank, existed = sl.GetRankByScore(2)
	assert.Equal(t, 0, rank)
	assert.Equal(t, true, existed)

	for n := sl.Front(); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}

	fmt.Println()
	n := sl.GetNodeLeastScore(-1)
	assert.Equal(t, int64(2), n.Score())
	assert.Equal(t, "b", n.Val())

	n = sl.GetNodeLeastScore(2)
	assert.Equal(t, int64(2), n.Score())
	assert.Equal(t, "b", n.Val())

	n = sl.GetNodeLeastScore(3)
	assert.Equal(t, int64(3), n.Score())
	assert.Equal(t, "c", n.Val())

	n = sl.GetNodeLeastScore(5)
	assert.Equal(t, int64(5), n.Score())
	assert.Equal(t, "a", n.Val())

	n = sl.GetNodeLeastScore(8)
	assert.Nil(t, n)
}
