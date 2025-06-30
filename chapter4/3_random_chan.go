package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"text/tabwriter"
)

func randomStream(stop <-chan struct{}, min int, max int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for {
			select {
			case <-stop:
				return
			default:
				res <- rand.IntN(max-min+1) + min
			}
		}
	}()
	return res
}

func main() {
	stop := make(chan struct{})
	minV, maxV := 1, 10
	randOneToTenCh := randomStream(stop, minV, maxV)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	tw.Write([]byte("N\tmin/max\trand\n"))
	for i := 0; i < 10; i++ {
		if v, ok := <-randOneToTenCh; ok {
			tw.Write([]byte(fmt.Sprintf("%d\t[%d, %d]\t%d\n", i+1, minV, maxV, v)))
		}
	}
	tw.Flush()
}
