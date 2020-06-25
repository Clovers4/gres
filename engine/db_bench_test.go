package engine

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkUnixtime(b *testing.B) {
	now := time.Now().Unix()
	for i := 0; i < b.N; i++ {
		time.Unix(now, 0)
	}
}

var db = NewDB()

func TestDB(t *testing.T) {
	concurrent := 10
	ops := 50000

	startG := sync.WaitGroup{}
	goG := sync.WaitGroup{}
	goG.Add(1)
	startG.Add(concurrent)
	endG := sync.WaitGroup{}
	endG.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go func() {
			startG.Done()
			goG.Wait()
			for i := 0; i < ops; i++ {
				key := strconv.FormatInt(int64(i), 10)
				val := key
				db.Set(key, val)
			}
			endG.Done()
		}()
	}
	startG.Wait()
	goG.Done()
	startTime := time.Now()
	endG.Wait()
	endTime := time.Now()

	use := endTime.Sub(startTime).Seconds()
	per := strconv.FormatFloat(float64(concurrent*ops)/float64(use), 'f', -1, 64)
	fmt.Println(use, per)
	fmt.Print(db.DbSize())
}

func TestDB2(t *testing.T) {
	concurrent := 10
	ops := 50000

	startG := sync.WaitGroup{}
	goG := sync.WaitGroup{}
	goG.Add(1)
	startG.Add(concurrent)
	endG := sync.WaitGroup{}
	endG.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go func() {
			startG.Done()
			goG.Wait()
			for i := 0; i < ops; i++ {
				key := strconv.FormatInt(int64(i), 10)
				val := key
				db.Set(key, val)
			}
			endG.Done()
		}()
	}
	startG.Wait()
	goG.Done()
	startTime := time.Now()
	endG.Wait()
	endTime := time.Now()

	use := endTime.Sub(startTime).Seconds()
	per := strconv.FormatFloat(float64(concurrent*ops)/float64(use), 'f', -1, 64)
	fmt.Println(use, per)
	fmt.Print(db.DbSize())
}

func BenchmarkDB(b *testing.B) {

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
	}
}
