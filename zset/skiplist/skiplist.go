package skiplist

import (
	"math/rand"
	"time"
)

const (
	MaxLevel    = 64
	Probability = 0.25
)

type SkiplistNode struct {
	score int64
	val   string

	prev      *SkiplistNode
	levels    []skiplistLevel
	probTable []float64
}

func newNode(level int, score int64, val string) *SkiplistNode {
	return &SkiplistNode{
		score:  score,
		val:    val,
		levels: make([]skiplistLevel, level, level),
	}
}

func (n *SkiplistNode) Next() *SkiplistNode {
	return n.levels[0].next
}

func (n *SkiplistNode) Prev() *SkiplistNode {
	return n.prev
}

func (n *SkiplistNode) Score() int64 {
	return n.score
}

func (n *SkiplistNode) Val() string {
	return n.val
}

type skiplistLevel struct {
	next *SkiplistNode
	span int
}

type Skiplist struct {
	header *SkiplistNode
	tail   *SkiplistNode
	level  int // 最大节点层数
	length int // 节点数量

	randSource rand.Source
}

func New() *Skiplist {
	return &Skiplist{
		level:  1,
		length: 0,
		header: newNode(MaxLevel, 0, ""),

		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (sl *Skiplist) Insert(score int64, val string) *SkiplistNode {
	update := make([]*SkiplistNode, MaxLevel, MaxLevel)
	rank := make([]int, MaxLevel, MaxLevel)

	curr := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// 直到下一个 score / val 大于 curr
		for next := curr.levels[i].next; next != nil && (next.score < score || next.score == score && next.val < val); next = curr.levels[i].next {

			rank[i] += curr.levels[i].span
			curr = curr.levels[i].next
		}

		update[i] = curr
	}

	// 假设节点不存在, 且我们允许重复 score , 并且有 map 进行去重
	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].levels[i].span = sl.length
		}
		sl.level = level
	}

	n := newNode(level, score, val)
	for i := 0; i < level; i++ {
		n.levels[i].next = update[i].levels[i].next
		update[i].levels[i].next = n

		n.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = (rank[0] - rank[i]) + 1
	}

	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	if update[0] != sl.header {
		n.prev = update[0]
	}

	if n.levels[0].next != nil {
		n.levels[0].next.prev = n
	} else {
		sl.tail = n
	}
	sl.length++
	return n
}

func (sl *Skiplist) Delete(score int64, val string) {
	update := make([]*SkiplistNode, MaxLevel, MaxLevel)

	n := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for n.levels[i].next != nil &&
			(n.levels[i].next.score < score ||
				(n.levels[i].next.score == score &&
					n.levels[i].next.val < val)) {
			n = n.levels[i].next
		}
		update[i] = n
	}

	n = n.levels[0].next
	// found
	if n != nil && n.score == score && n.val == val {
		sl.deleteNode(n, update)
	}
	// not found
}

func (sl *Skiplist) deleteNode(n *SkiplistNode, update []*SkiplistNode) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].next == n {
			update[i].levels[i].span += n.levels[i].span - 1
			update[i].levels[i].next = n.levels[i].next
		} else {
			update[i].levels[i].span--
		}
	}
	if n.levels[0].next == nil {
		sl.tail = n.prev
	} else {
		n.levels[0].next.prev = n.prev
	}

	for sl.level > 1 && sl.header.levels[sl.level-1].next == nil {
		sl.level--
	}
	sl.length--
}

// 注意 val==n.val
func (sl *Skiplist) UpdateScore(curScore, newScore int64, val string) *SkiplistNode {
	update := make([]*SkiplistNode, MaxLevel, MaxLevel)

	n := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for n.levels[i].next != nil &&
			(n.levels[i].next.score < curScore ||
				(n.levels[i].next.score == curScore &&
					n.levels[i].next.val < val)) {
			n = n.levels[i].next
		}
		update[i] = n
	}

	n = n.levels[0].next
	if (n.prev == nil || n.prev.score < newScore) &&
		(n.levels[0].next == nil || n.levels[0].next.score > newScore) {
		n.score = newScore
	}

	// cannot reuse old node
	sl.deleteNode(n, update)
	newNode := sl.Insert(newScore, n.val)
	return newNode
}

func (sl *Skiplist) Front() *SkiplistNode {
	return sl.header.levels[0].next
}

func (sl *Skiplist) End() *SkiplistNode {
	return sl.tail
}

// rank start at 1, end at sl.length
func (sl *Skiplist) GetNodeByRank(rank int) *SkiplistNode {
	if rank < 0 {
		rank = sl.length + rank
	}
	rank++
	if rank <= 0 || rank > sl.length {
		return nil
	}

	traversed := 0
	n := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for n.levels[i].next != nil && (traversed+n.levels[i].span) <= rank {
			traversed += n.levels[i].span
			n = n.levels[i].next
		}
		if traversed == rank {
			return n
		}
	}
	return nil
}

// the rank start at 0, end at length-1
func (sl *Skiplist) GetRankByScore(score int64) (rank int, existed bool) {
	traversed := 0
	n := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for n.levels[i].next != nil && n.levels[i].next.score <= score {
			traversed += n.levels[i].span
			n = n.levels[i].next
		}
		if n.score == score {
			return traversed - 1, true
		}
	}
	return -1, false
}

func (sl *Skiplist) GetNodeLeastScore(score int64) *SkiplistNode {
	n := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for n.levels[i].next != nil && n.levels[i].next.score < score {
			n = n.levels[i].next
		}
	}
	return n.levels[0].next
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for sl.random() < Probability {
		level++
	}
	if level > MaxLevel {
		return MaxLevel
	}
	return level
}

func (sl *Skiplist) random() float64 {
	return float64(sl.randSource.Int63()) / (1 << 63)
}

func (sl *Skiplist) Length() int {
	return sl.length
}
