package main

import (
	"sync"
	"sync/atomic"
	"time"
)

const concurrency = 10_000_000

var counter int64

func process() {
	var wg sync.WaitGroup
	wg.Add(concurrency)

	chLimit := make(chan struct{}, 500_000)

	for i := 0; i < concurrency; i++ {
		chLimit <- struct{}{}
		go func() {
			defer func() {
				<-chLimit
				wg.Done()
			}()
			atomic.AddInt64(&counter, 1)
			time.Sleep(3 * time.Second)
		}()
	}
	wg.Wait()
}
