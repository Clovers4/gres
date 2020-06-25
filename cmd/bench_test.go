package cmd

import (
	"fmt"
	"github.com/clovers4/gres"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestDB2(t *testing.T) {
	// init and start server
	go func() {
		srv := gres.NewServer()
		defer srv.Stop()
		srv.Start()
	}()

	concurrent := 100
	ops := 500

	startG := sync.WaitGroup{}
	goG := sync.WaitGroup{}
	goG.Add(1)
	startG.Add(concurrent)
	endG := sync.WaitGroup{}
	endG.Add(concurrent)

	cli := NewClient()
	cli.InitSet(ops)
	for i := 0; i < concurrent; i++ {
		go func() {
			cli := NewClient()
			startG.Done()
			goG.Wait()
			cli.BenchSet(ops)
			endG.Done()
			cli.GracefulExit()
		}()
	}
	startG.Wait()
	goG.Done()
	startTime := time.Now()
	endG.Wait()
	endTime := time.Now()

	use := endTime.Sub(startTime).Seconds()
	per := strconv.FormatFloat(float64(concurrent*ops)/float64(use), 'f', -1, 64)
	fmt.Println(fmt.Sprintf("共使用%v秒，每秒将执行%v次Set操作", use, per))
}
