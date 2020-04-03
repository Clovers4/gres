package skiplist

import (
	"github.com/clovers4/gres/engine"
	"math/rand"
	"strconv"
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
)

func BenchmarkSkiplistSequenceInsert(b *testing.B) {
	sl := engine.New()
	for i := 0; i < b.N; i++ {
		sl.Insert(float64(i), strconv.FormatInt(int64(i), 10))
	}
}

func BenchmarkSkiplistRandomInsert(b *testing.B) {
	sl := engine.New()
	for i := 0; i < b.N; i++ {
		sl.Insert(float64(i%10), strconv.FormatInt(int64(i), 10))
	}
}

func BenchmarkArraylistSequenceInsert(b *testing.B) {
	ls := arraylist.New()
	for i := 0; i < b.N; i++ {
		ls.Insert(ls.Size(), float64(i))
	}
}

func BenchmarkArraylistRandomInsert(b *testing.B) {
	ls := arraylist.New()
	ls.Insert(0, 1)
	b.N = 100000

	for i := 0; i < b.N; i++ {
		ls.Insert(rand.Int()%ls.Size(), float64(i))
	}
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}
