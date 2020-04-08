package engine

import (
	"testing"
	"time"
)

func BenchmarkUnixtime(b *testing.B) {
	now := time.Now().Unix()
	for i := 0; i < b.N; i++ {
		time.Unix(now, 0)
	}
}
