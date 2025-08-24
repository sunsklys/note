package g_limit

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"
)

func TestCh(t *testing.T) {
	limit := NewCh(10)
	limit.Start(math.MaxInt8)
	limit.wg.Wait()
}

type Ch struct {
	wg    *sync.WaitGroup
	ch    chan struct{}
	limit int
}

func NewCh(limit int) *Ch {
	return &Ch{
		wg:    &sync.WaitGroup{},
		ch:    make(chan struct{}, limit),
		limit: limit,
	}
}

func (r *Ch) Start(n int) {
	for i := 0; i < n; i++ {
		r.wg.Add(1)
		r.ch <- struct{}{}
		go r.Deal(i)
	}
	close(r.ch)
}

func (r *Ch) Deal(i int) {
	defer func() {
		<-r.ch
		r.wg.Done()
	}()
	fmt.Println(i, runtime.NumGoroutine())
}
