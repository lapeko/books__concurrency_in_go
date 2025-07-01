package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func tee(in <-chan int) (<-chan int, <-chan int) {
	out1, out2 := make(chan int), make(chan int)

	go func() {
		defer func() {
			close(out1)
			close(out2)
		}()

		for num := range in {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case out1 <- num:
					out1 = nil
				case out2 <- num:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func generateNInts(size int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := 1; num <= size; num++ {
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(500)))
			out <- num
		}
	}()

	return out
}

func main() {
	size := 10
	out1, out2 := tee(generateNInts(size))
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range out1 {
			fmt.Println("out1:", num)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range out2 {
			fmt.Println("out2:", num)
		}
	}()

	wg.Wait()
}
