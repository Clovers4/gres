package skiplist

import (
	"fmt"
	"github.com/clovers4/gres/engine"
	"testing"
)

func TestSkiplist(t *testing.T) {
	sl := engine.New()
	sl.Insert(1.0, "a")
	sl.Insert(1.0, "a22")
	sl.Insert(1.0, "a")
	sl.Insert(3.0, "c")
	sl.Insert(2.0, "b")

	for n := sl.Front(); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	sl.Delete(1.0, "a22")

	for n := sl.Front(); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	sl.UpdateScore(1.0, 5, "a")

	for n := sl.Front(); n != nil; n = n.Next() {
		fmt.Println(n.Score(), n.Val())
	}
	fmt.Println()

	target := sl.GetNodeByRank(1)
	fmt.Println(target.Score(), target.Val())

	target = sl.GetNodeByRank(2)
	fmt.Println(target.Score(), target.Val())

	target = sl.GetNodeByRank(0)
	if target != nil {
		fmt.Println(target.Score(), target.Val())
	}

	target = sl.GetNodeByRank(5)
	if target != nil {
		fmt.Println(target.Score(), target.Val())
	}
}
