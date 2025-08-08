package paradigm

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"
)

func TestGLimitConsumer(t *testing.T) {
	limit := NewGLimitConsumer(10)
	limit.Consumer()
	limit.Producer(math.MaxInt8)
	limit.wg.Wait()
}

type GLimitConsumer struct {
	wg    *sync.WaitGroup
	ch    chan int
	limit int
}

func NewGLimitConsumer(limit int) *GLimitConsumer {
	return &GLimitConsumer{
		wg:    &sync.WaitGroup{},
		ch:    make(chan int, limit),
		limit: limit,
	}
}

func (r *GLimitConsumer) Producer(n int) {
	for i := 0; i < n; i++ {
		r.wg.Add(1)
		r.ch <- i
	}
	close(r.ch)
}

func (r *GLimitConsumer) Consumer() {
	for i := 0; i < r.limit; i++ {
		go func(j int) {
			for v := range r.ch {
				r.wg.Done()
				fmt.Println(v, j, runtime.NumGoroutine())
			}
		}(i)
	}
}
