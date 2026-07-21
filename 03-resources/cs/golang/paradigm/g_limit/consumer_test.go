package g_limit

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"
)

func TestConsumer(t *testing.T) {
	limit := NewConsumer(10)
	limit.Consumer()
	limit.Producer(math.MaxInt8)
	limit.wg.Wait()
}

type Consumer struct {
	wg    *sync.WaitGroup
	ch    chan int
	limit int
}

func NewConsumer(limit int) *Consumer {
	return &Consumer{
		wg:    &sync.WaitGroup{},
		ch:    make(chan int, limit),
		limit: limit,
	}
}

func (r *Consumer) Producer(n int) {
	for i := 0; i < n; i++ {
		r.wg.Add(1)
		r.ch <- i
	}
	close(r.ch)
}

func (r *Consumer) Consumer() {
	for i := 0; i < r.limit; i++ {
		go func(j int) {
			for v := range r.ch {
				r.wg.Done()
				fmt.Println(v, j, runtime.NumGoroutine())
			}
		}(i)
	}
}
