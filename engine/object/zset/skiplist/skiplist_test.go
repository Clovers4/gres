package skiplist

import (
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
	sl.Insert(1.1, "a")
	sl.Insert(1.1, "a22")
	sl.Insert(3.1, "c")
	sl.Insert(2.1, "b")

	sl.Delete(1.1, "a22")

	var vals []interface{}
	for n := sl.Front(); n != nil; n = n.Next() {
		vals = append(vals, n.Score(), n.Val())
	}
	assert.Equal(t, 6, len(vals))
	assert.Equal(t, "{1.1, a, 2.1, b, 3.1, c}", array2String(vals))

	sl.UpdateScore(1.1, 5.1, "a")
	vals = []interface{}{}
	for n := sl.Front(); n != nil; n = n.Next() {
		vals = append(vals, n.Score(), n.Val())
	}
	assert.Equal(t, 6, len(vals))
	assert.Equal(t, "{2.1, b, 3.1, c, 5.1, a}", array2String(vals))

	target := sl.GetNodeByRank(0)
	assert.Equal(t, 2.1, target.Score())
	assert.Equal(t, "b", target.Val())

	target = sl.GetNodeByRank(1)
	assert.Equal(t, 3.1, target.Score())
	assert.Equal(t, "c", target.Val())

	target = sl.GetNodeByRank(-1)
	assert.Nil(t, target)

	target = sl.GetNodeByRank(5)
	assert.Nil(t, target)

	var rank int
	var existed bool
	rank, existed = sl.GetRankByScore(5.5)
	assert.Equal(t, false, existed)

	rank, existed = sl.GetRankByScore(5.1)
	assert.Equal(t, 2, rank)
	assert.Equal(t, true, existed)

	rank, existed = sl.GetRankByScore(2.1)
	assert.Equal(t, 0, rank)
	assert.Equal(t, true, existed)
}
