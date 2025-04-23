package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var cadence = sync.NewCond(&sync.Mutex{})

func takeStep() {
	cadence.L.Lock()
	cadence.Wait()
	cadence.L.Unlock()
}

func tryDir(dirName string, dir *int32, out *bytes.Buffer) bool {
	fmt.Fprintf(out, " %v", dirName)
	atomic.AddInt32(dir, 1)
	takeStep()
	if atomic.LoadInt32(dir) == 1 {
		fmt.Fprint(out, ". Success!")
		return true
	}
	fmt.Fprint(out, ". Failed!\n")
	takeStep()
	atomic.AddInt32(dir, -1)
	return false
}

var left, right int32

func tryLeft(out *bytes.Buffer) bool {
	return tryDir("left", &left, out)
}

func tryRight(out *bytes.Buffer) bool {
	return tryDir("right", &right, out)
}

func walk(wg *sync.WaitGroup, name string) {
	var out bytes.Buffer
	defer func() {
		fmt.Println(out.String())
	}()
	defer wg.Done()
	fmt.Fprintf(&out, "%v is trying to scoot:", name)
	for i := 0; i < 3; i++ {
		if tryLeft(&out) || tryRight(&out) {
			return
		}
	}
	fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
}

func main() {
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
