package main

import (
	"bytes"
	"fmt"
	"time"
)

func collectLogs(
	logStream <-chan string,
	terminate <-chan struct{},
) <-chan string {
	res := make(chan string)
	var buff bytes.Buffer

	go func() {
		defer close(res)
		for {
			select {
			case log := <-logStream:
				buff.Write([]byte(log + "\n"))
			case <-terminate:
				res <- string(buff.Bytes())
				buff.Reset()
				return
			}
		}
	}()

	return res
}

func collectLogsSync(
	logStream <-chan string,
	terminate <-chan struct{},
) string {
	return <-collectLogs(logStream, terminate)
}

func main() {
	logStream := make(chan string)
	terminate := make(chan struct{})

	go func() {
		logStream <- "[GET]: Success"
		logStream <- "[POST]: Success"
		logStream <- "[GET]: Failure. Not found"
		logStream <- "[POST]: Success"
		time.Sleep(time.Second)
		terminate <- struct{}{}
	}()

	allLogs := collectLogsSync(logStream, terminate)

	fmt.Println(allLogs)
}
