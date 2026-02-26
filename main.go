package main

import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	Id  int
	Url string
}

func worker(id int, in <-chan Server, out chan<- string) {
	for server := range in {
		fmt.Printf("worker %d check url %s\n", id, server.Url)
		time.Sleep(3 * time.Second)
		out <- fmt.Sprintf("url %s checked at %s", server.Url, time.Now())
	}
}

func producer(in chan<- Server, qty int) {
	for i := range qty {
		server := Server{Id: i, Url: fmt.Sprintf("https://site%d.com", i)}
		in <- server
	}
	close(in)
}

func consumer(out <-chan string, done chan struct{}) {
	for i := range out {
		fmt.Println(i)
	}
	close(done)
}

func checkServers() {
	qty := 5
	workersCont := 3
	var wg sync.WaitGroup
	in := make(chan Server, qty)
	out := make(chan string, qty)
	done := make(chan struct{})

	for i := range workersCont {
		wg.Go(func() {
			worker(i, in, out)
		})
	}

	go func() {
		producer(in, qty)
	}()

	go func() {
		consumer(out, done)
	}()

	wg.Wait()
	close(out)

	<-done
	fmt.Println("all workers done")
}

func main() {
	ticker := time.NewTicker(time.Second * 20)

	for range ticker.C {
		checkServers()
	}
}
