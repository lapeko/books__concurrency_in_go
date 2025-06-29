package main

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

type HeavyDecoder struct {
	memory []byte
}

func (decoder *HeavyDecoder) decode() {}

var heavyDecoderPool = func(size int) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &HeavyDecoder{memory: make([]byte, 0, size)}
		},
	}
}

func main() {
	iterations := int(1e4)
	size := int(1e6)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hd := &HeavyDecoder{memory: make([]byte, size)}
			hd.decode()
		}()
	}
	wg.Wait()

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintf(tw, "%s\t%s\n", "Pool used", "Time takes")
	fmt.Fprintf(tw, "%t\t%s\n", false, time.Since(start))

	hdp := heavyDecoderPool(size)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hd := hdp.Get().(*HeavyDecoder)
			hdp.Put(hd)
			hd.decode()
		}()
	}
	wg.Wait()
	fmt.Fprintf(tw, "%t\t%s\n", true, time.Since(start))
	tw.Flush()
}
