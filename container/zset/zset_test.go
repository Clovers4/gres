package zset

import (
	"fmt"
	"testing"
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
