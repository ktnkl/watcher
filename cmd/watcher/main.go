package main

import (
	"fmt"
	"sync"
	"time"
	"watcher/internal/consumer"
	"watcher/internal/domain"
	"watcher/internal/producer"
	"watcher/internal/worker"
)

func checkServers() {
	qty := 5
	workersCont := 3
	var wg sync.WaitGroup
	in := make(chan domain.Server, qty)
	out := make(chan string, qty)
	done := make(chan struct{})

	timer := worker.RealTimer{}
	logger := worker.RealLogger{}

	for i := range workersCont {
		wg.Go(func() {
			worker.Worker(i, in, out, timer, logger)
		})
	}

	go func() {
		producer.Producer(in, qty)
	}()

	go func() {
		consumer.Consumer(out, done)
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
