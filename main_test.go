package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func CreateTestServers(count int) []Server {
	servers := make([]Server, count)

	for c := range count {
		servers[c] = Server{
			Id:  c,
			Url: fmt.Sprintf("https://test%d.com", c),
		}
	}

	return servers
}

func TestWorker(t *testing.T) {
	in := make(chan Server, 1)
	out := make(chan string, 1)
	timer := &MockTimer{NowTime: time.Now()}
	logger := &MockLogger{}

	go worker(1, in, out, timer, logger)

	in <- Server{Id: 1, Url: "https://test1.com"}
	close(in)

	result := <-out

	if !timer.SleepCalled {
		t.Error("worker не вызвал sleep")
	}

	if len(logger.Messages) == 0 {
		t.Error("worker не логгировал")
	}

	if result == "" {
		t.Error("пустой результат")
	}
}

func TestWorkerTableDriven(t *testing.T) {
	var tests = []struct {
		qty       int
		servers   []Server
		wantedUrl []string
	}{
		{
			1,
			CreateTestServers(1),
			[]string{"https://test0.com"},
		},
		{
			3,
			CreateTestServers(3),
			[]string{"https://test0.com", "https://test1.com", "https://test2.com"},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("worker table test, servers qty: %d", tt.qty)

		t.Run(testname, func(t *testing.T) {
			var wg sync.WaitGroup
			in := make(chan Server, tt.qty)
			out := make(chan string, tt.qty)

			timer := &MockTimer{NowTime: time.Now()}
			logger := &MockLogger{}

			wg.Go(func() {
				worker(1, in, out, timer, logger)
			})

			for _, server := range tt.servers {
				in <- server
			}
			close(in)

			wg.Wait()
			close(out)

			result := make([]string, tt.qty)

			for i := range tt.qty {
				server := <-out
				result[i] = server
			}

			if !timer.SleepCalled {
				t.Error("worker не вызвал sleep")
			}

			if len(logger.Messages) != tt.qty {
				t.Error("worker не логгировал")
			}

			if len(result) != len(tt.wantedUrl) {
				t.Error("Кол-во ожидаемых и полученных результатов не совпадает")
			}

			for i, result := range result {
				if !strings.Contains(result, tt.wantedUrl[i]) {
					t.Error("Ожидаемые ответы не сопадают с реальными")
				}
			}

		})
	}

}
