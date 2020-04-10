package main

import (
	"fmt"
	"github.com/clovers4/gres"
	"sync"
	"testing"
	"time"
)

func TestBenchmark(t *testing.T) {
	go func() {
		srv := gres.NewServer()
		srv.Start()
		defer srv.Stop()
	}()
	time.Sleep(2 * time.Second)

	threadNum := 20
	opNum := 2000

	startWg := sync.WaitGroup{}
	startWg.Add(1)
	endWg := sync.WaitGroup{}
	endWg.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		go func() {
			cli := NewClient()
			startWg.Wait()
			for n := 0; n < opNum; n++ {
				cli.interact(fmt.Sprintf("set key%v value%v0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", n+i*threadNum, n))
			}
			endWg.Done()
		}()
	}
	startTime := time.Now()
	startWg.Done()
	endWg.Wait()
	uesd := time.Now().Sub(startTime)
	fmt.Println(threadNum, opNum, threadNum*opNum)
	fmt.Println("use:", uesd, float64(threadNum*opNum)/uesd.Seconds())

}

func TestBenchmarkGet(t *testing.T) {
	go func() {
		srv := gres.NewServer(gres.DbnumOption(8))
		srv.Start()
		defer srv.Stop()
	}()
	time.Sleep(2 * time.Second)

	threadNum := 20
	opNum := 2000

	startWg := sync.WaitGroup{}
	startWg.Add(1)
	endWg := sync.WaitGroup{}
	endWg.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		go func() {
			cli := NewClient()
			startWg.Wait()
			for n := 0; n < opNum; n++ {
				cli.interact(fmt.Sprintf("get key%v", n+i*threadNum))
			}
			endWg.Done()
		}()
	}
	startTime := time.Now()
	startWg.Done()
	endWg.Wait()
	uesd := time.Now().Sub(startTime)
	fmt.Println(threadNum, opNum, threadNum*opNum)
	fmt.Println("use:", uesd, float64(threadNum*opNum)/uesd.Seconds())

}
