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

type Timer interface {
	Sleep(time.Duration)
	Now() time.Time
}

type Logger interface {
	Printf(format string, args ...interface{})
}

type RealTimer struct{}

func (RealTimer) Now() time.Time {
	return time.Now()
}

func (RealTimer) Sleep(d time.Duration) {
	time.Sleep(d)
}

type RealLogger struct{}

func (RealLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

type MockTimer struct {
	NowTime       time.Time
	SleepCalled   bool
	SleepDuration time.Duration
}

func (m *MockTimer) Now() time.Time {
	return m.NowTime
}

func (m *MockTimer) Sleep(d time.Duration) {
	m.SleepCalled = true
	m.SleepDuration = d
}

type MockLogger struct {
	Messages []string
}

func (m *MockLogger) Printf(format string, args ...interface{}) {
	m.Messages = append(m.Messages, fmt.Sprintf(format, args...))
}

func worker(id int, in <-chan Server, out chan<- string, timer Timer, logger Logger) {
	for server := range in {
		logger.Printf("worker %d check url %s\n", id, server.Url)
		timer.Sleep(3 * time.Second)
		out <- fmt.Sprintf("url %s checked at %s", server.Url, timer.Now())
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

	timer := RealTimer{}
	logger := RealLogger{}

	for i := range workersCont {
		wg.Go(func() {
			worker(i, in, out, timer, logger)
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
