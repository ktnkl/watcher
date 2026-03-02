package consumer

import "fmt"

func Consumer(out <-chan string, done chan struct{}) {
	for i := range out {
		fmt.Println(i)
	}
	close(done)
}
