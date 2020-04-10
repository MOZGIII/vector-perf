package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var totalBytesRead uint64
var totalRequests uint64

func handle(w http.ResponseWriter, req *http.Request) {
	bytesRead, err := io.Copy(ioutil.Discard, req.Body)
	if err != nil {
		panic(err)
	}
	atomic.AddUint64(&totalRequests, 1)
	atomic.AddUint64(&totalBytesRead, uint64(bytesRead))
	time.Sleep(100 * time.Millisecond)
	w.WriteHeader(200)
}

func printStats() {
	for {
		bytesSecondAgo := atomic.LoadUint64(&totalBytesRead)
		requestsSecondAgo := atomic.LoadUint64(&totalRequests)
		time.Sleep(3 * time.Second)
		bytesNow := atomic.LoadUint64(&totalBytesRead)
		requestsNow := atomic.LoadUint64(&totalRequests)
		fmt.Printf(
			"B: %10d (%14d) R: %5d (%10d)\n",
			bytesNow-bytesSecondAgo, bytesNow,
			requestsNow-requestsSecondAgo, requestsNow,
		)
	}
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":9001"
	}
	fmt.Printf("Listening at %s\n", addr)
	go printStats()
	if err := http.ListenAndServe(addr, http.HandlerFunc(handle)); err != nil {
		panic(err)
	}
}
