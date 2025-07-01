package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func buildStream(list []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, num := range list {
			//time.Sleep(time.Millisecond * time.Duration(rand.IntN(1000)))
			out <- num
		}
		close(out)
	}()
	return out
}

func heavyCalc(num int) int {
	time.Sleep(time.Second)
	return num + 1
}

func runConsistently(in <-chan int) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case num, notClosed := <-in:
				if !notClosed {
					return
				}
				fmt.Printf("Consistent: %d\n", heavyCalc(num))
			}
		}
	}()

	wg.Wait()
	fmt.Printf("Run consistently time: %s\n", time.Since(start))
}

func runConcurrently(in <-chan int) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for num := range runAsFanOutIn(in, heavyCalc, heavyCalc, heavyCalc, heavyCalc) {
			fmt.Printf("Parallel: %d\n", num)
		}
	}()

	wg.Wait()
	fmt.Printf("Run in parallel time: %s\n", time.Since(start))
}

func runAsFanOutIn(in <-chan int, handlers ...func(int) int) <-chan int {
	out := make(chan int)

	var notCompletedHandlers atomic.Int64
	notCompletedHandlers.Store(int64(len(handlers)))

	for _, handler := range handlers {
		go func() {
			for {
				select {
				case num, notCompleted := <-in:
					if notCompleted {
						out <- handler(num)
					} else if notCompletedHandlers.Load() == 1 {
						close(out)
						return
					} else {
						notCompletedHandlers.Add(-1)
						return
					}
				}
			}
		}()
	}

	return out
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	go func() {
		defer wg.Done()
		runConsistently(buildStream(in))
	}()

	go func() {
		defer wg.Done()
		runConcurrently(buildStream(in))
	}()

	wg.Wait()
}
