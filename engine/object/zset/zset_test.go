package zset

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZSet(t *testing.T) {
	zs := New()
	zs.Add(1.0, "A")
	zs.Add(1.0, "A")
	zs.Add(2.0, "B")
	zs.Add(2.0, "B1")
	zs.Add(3.0, "C")

	assert.Equal(t, true, zs.Length() == 4)

	for n := zs.GetNodeByRank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	zs.Add(3.0, "A")

	for n := zs.GetNodeByRank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	zs.Delete("A")
	zs.Delete("Unknwon")
	assert.Equal(t, true, zs.Length() == 3)
	_, ok := zs.Get("A")
	assert.Equal(t, false, ok)

	for n := zs.GetNodeByRank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
}

func TestZSetMarshal(t *testing.T) {
	var err error

	zs := New()
	zs.Add(1.32, "A")
	zs.Add(2.23, "B1")
	zs.Add(3.44, "C")

	// marshal
	buf := new(bytes.Buffer)
	err = zs.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newZs := New()
	r := bytes.NewReader(buf.Bytes())
	err = newZs.Unmarshal(r)
	assert.Nil(t, err)

	assert.Equal(t, 3, newZs.Length())
	for i := 1; i <= newZs.Length(); i++ {
		old, new := zs.GetNodeByRank(i), newZs.GetNodeByRank(i)
		assert.Equal(t, old.Val(), new.Val())
		assert.Equal(t, old.Score(), new.Score())
	}
}
