package worker

import (
	"fmt"
	"time"
	"watcher/internal/domain"
)

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

func Worker(id int, in <-chan domain.Server, out chan<- string, timer Timer, logger Logger) {
	for server := range in {
		logger.Printf("worker %d check url %s\n", id, server.Url)
		timer.Sleep(3 * time.Second)
		out <- fmt.Sprintf("url %s checked at %s", server.Url, timer.Now())
	}
}
