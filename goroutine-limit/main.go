package main

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			pringMemoryUsage()
		}
	}()
	process()
}

func pringMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	w := tabwriter.NewWriter(os.Stdout, 18, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Alloc = %v MiB\t", convertToMB(m.Alloc))
	fmt.Fprintf(w, "NumGC = %v\t", m.NumGC)
	fmt.Fprintf(w, "Completed = %v%%\n", (atomic.LoadInt64(&counter)*100)/concurrency)
	w.Flush()
}

func convertToMB(b uint64) uint64 {
	return b / 1024 / 1024
}
