package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func writer(wg *sync.WaitGroup, l sync.Locker) {
	defer wg.Done()
	defer l.Unlock()
	l.Lock()
	time.Sleep(time.Millisecond * 50)
}

func reader(wg *sync.WaitGroup, l sync.Locker) {
	defer wg.Done()
	defer l.Unlock()
	l.Lock()
	time.Sleep(time.Millisecond * 10)
}

func test(count int, mutex, rwMutex sync.Locker) time.Duration {
	var wg sync.WaitGroup
	wg.Add(count + 1)
	beginTestTime := time.Now()
	go writer(&wg, mutex)
	for i := count; i > 0; i-- {
		go reader(&wg, rwMutex)
	}
	wg.Wait()
	return time.Since(beginTestTime)
}

func main() {
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()
	fmt.Fprintf(tw, "Readers\tRWMutext\tMutex\n")

	var m sync.RWMutex
	for i := 0; i < 10; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}
