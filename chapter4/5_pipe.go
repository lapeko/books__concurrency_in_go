package main

import "fmt"

func makeIntChanList(ints []int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range ints {
			ch <- v
		}
	}()
	return ch
}

func add(incNum int) func(<-chan int) <-chan int {
	return func(in <-chan int) <-chan int {
		out := make(chan int)

		go func() {
			defer close(out)
			for {
				select {
				case currentNum, ok := <-in:
					if ok {
						out <- currentNum + incNum
					} else {
						return
					}
				}
			}
		}()

		return out
	}
}

func pipeMap[T any](mapFunc func(T) T) func(<-chan T) <-chan T {
	return func(in <-chan T) <-chan T {
		out := make(chan T)

		go func() {
			defer close(out)
			for {
				select {
				case current, ok := <-in:
					if ok {
						out <- mapFunc(current)
					} else {
						return
					}
				}
			}
		}()

		return out
	}
}

func pipe[T any](dataStream <-chan T, funcs ...func(<-chan T) <-chan T) <-chan T {
	for _, f := range funcs {
		dataStream = f(dataStream)
	}
	return dataStream
}

func main() {
	ch := pipe(
		makeIntChanList([]int{1, 2, 3, 4, 5}),
		add(3),
		pipeMap(func(x int) int { return x * 2 }),
	)

	for n := range ch {
		fmt.Println(n)
	}
}
