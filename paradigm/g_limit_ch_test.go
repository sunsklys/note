package paradigm

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"
)

func TestGLimitCh(t *testing.T) {
	limit := NewGLimit(10)
	limit.Start(math.MaxInt8)
	limit.wg.Wait()
}

type GLimitCh struct {
	wg    *sync.WaitGroup
	ch    chan struct{}
	limit int
}

func NewGLimit(limit int) *GLimitCh {
	return &GLimitCh{
		wg:    &sync.WaitGroup{},
		ch:    make(chan struct{}, limit),
		limit: limit,
	}
}

func (r *GLimitCh) Start(n int) {
	for i := 0; i < n; i++ {
		r.wg.Add(1)
		r.ch <- struct{}{}
		go r.Deal(i)
	}
	close(r.ch)
}

func (r *GLimitCh) Deal(i int) {
	defer func() {
		<-r.ch
		r.wg.Done()
	}()
	fmt.Println(i, runtime.NumGoroutine())
}
