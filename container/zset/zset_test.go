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

	for n := zs.Rank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	zs.Add(3.0, "A")

	for n := zs.Rank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	zs.Delete("A")
	zs.Delete("Unknwon")

	for n := zs.Rank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
}

func TestZSetMarshal(t *testing.T) {
	var err error

	zs := New()
	zs.Add(1.0, "A")
	zs.Add(2.0, "B1")
	zs.Add(3.0, "C")

	for n := zs.Rank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	buf := new(bytes.Buffer)
	err = zs.Marshal(buf)
	assert.Nil(t, err)
	fmt.Println(buf.Bytes())

	newZs := New()
	r := bytes.NewReader(buf.Bytes())
	err = newZs.Unmarshal(r)
	assert.Nil(t, err)
	fmt.Println(newZs.Length())
	for n := newZs.Rank(1); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
}
