package paradigm

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestWorkPool(t *testing.T) {
	wm := NewWorkerManager(10)
	// 定时触发退出信号，避免测试悬挂
	time.AfterFunc(2*time.Second, func() {
		wm.sig <- syscall.SIGINT
	})
	wm.StartWorkerPool()
}

type WorkerManager struct {
	workerChan chan *worker
	sig        chan os.Signal
	exit       bool
	nWorkers   int
	wg         *sync.WaitGroup
	lock       *sync.RWMutex
}

func NewWorkerManager(workers int) *WorkerManager {
	return &WorkerManager{
		nWorkers:   workers,
		workerChan: make(chan *worker, workers),
		sig:        make(chan os.Signal, 1),
		wg:         &sync.WaitGroup{},
		exit:       false,
		lock:       &sync.RWMutex{},
	}
}

func (r *WorkerManager) StartWorkerPool() {
	signal.Notify(r.sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-r.sig
		r.lock.Lock()
		r.exit = true
		r.lock.Unlock()
		close(r.workerChan)
	}()

	for i := 0; i < r.nWorkers; i++ {
		wk := &worker{id: i}
		r.wg.Add(1)
		go wk.work(r)
	}

	r.KeepLiveWorkers()
}

func (r *WorkerManager) KeepLiveWorkers() {
	for wk := range r.workerChan {
		r.wg.Add(1)
		fmt.Printf("Worker %d stopped with err: [%v] \n", wk.id, wk.err)
		wk.err = nil
		go wk.work(r)
	}
	r.wg.Wait()
}

type worker struct {
	id  int
	err error
}

func (w *worker) work(manager *WorkerManager) {
	var err error
	defer func() {
		manager.wg.Done()
		if r := recover(); r != nil {
			w.err = fmt.Errorf("panic happened with [%v]", r)
		} else {
			w.err = err
		}

		fmt.Println("Stop Worker...ID = ", w.id)
		manager.lock.RLock()
		if !manager.exit {
			manager.workerChan <- w
		}
		manager.lock.RUnlock()
	}()

	fmt.Println("Start Worker...ID = ", w.id)

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
	}

	if rand.Intn(10) > 5 {
		panic("worker panic..")
	} else {
		err = errors.New("work exit")
		runtime.Goexit()
	}
}
