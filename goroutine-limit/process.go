package main

import (
	"sync"
	"sync/atomic"
	"time"
)

const concurrency = 1_000_000

var counter int64

func process() {
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
			time.Sleep(3 * time.Second)
		}()
	}
	wg.Wait()
}
