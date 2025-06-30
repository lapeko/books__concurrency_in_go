package main

import "fmt"

func logFromChannel[T any](ch chan T) {
	value, ok := <-ch
	fmt.Printf("value: %v. ok: %t\n", value, ok)
}

func main() {
	ch := make(chan int)
	close(ch)

	logFromChannel(ch)

	ch = make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	logFromChannel(ch)
	for n := range ch {
		fmt.Println(n)
	}
	logFromChannel(ch)
}
