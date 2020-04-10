package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func handle(w http.ResponseWriter, req *http.Request) {
	_, err := io.Copy(ioutil.Discard, req.Body)
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
	w.WriteHeader(200)
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":9001"
	}
	fmt.Printf("Listening at %s\n", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(handle)); err != nil {
		panic(err)
	}
}
