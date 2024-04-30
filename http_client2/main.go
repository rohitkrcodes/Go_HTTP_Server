package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var ch = make(chan int)

type sku struct {
	item, price string
}

var items = []sku{
	{"shoes", "46"},
	{"socks", "16"},
	{"pants", "50"},
	{"shorts", "96"},
}

func nums(ch chan<- int) {
	for i := 0; ; i++ {
		ch <- i
	}
}

func doQuery(cmd, param string) error {
	resp, err := http.Get("http://localhost:8080/" + cmd + "?" + param)

	if err != nil {
		fmt.Fprintf(os.Stderr, "got %s=%d ERROR\n", param, err)
		return err
	}

	defer resp.Body.Close()

	x := <-ch
	fmt.Printf("number of requests processed %d\n", x)
	fmt.Fprintf(os.Stderr, "got %s=%d no error\n", param, resp.StatusCode)
	return nil
}

func runAds() {
	for {
		for _, s := range items {
			if err := doQuery("create", "item="+s.item+"&price="+s.price); err != nil {
				return
			}
		}
	}
}

func runUpdates() {
	for {
		for _, s := range items {
			if err := doQuery("update", "item="+s.item+"&price="+s.price); err != nil {
				return
			}
		}
	}
}

func runDeletes() {
	for {
		for _, s := range items {
			if err := doQuery("update", "item="+s.item); err != nil {
				return
			}
		}
	}
}

func main() {
	go runAds()
	go runUpdates()
	go runDeletes()
	go nums(ch)

	time.Sleep(5 * time.Second)
}
