package producer

import (
	"fmt"
	"watcher/internal/domain"
)

func Producer(in chan<- domain.Server, qty int) {
	for i := range qty {
		server := domain.Server{Id: i, Url: fmt.Sprintf("https://site%d.com", i)}
		in <- server
	}
	close(in)
}
