package wait

import (
	"sync"
	"time"
)

type Wait struct {
	wg sync.WaitGroup
}

func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

func (w *Wait) Done() {
	w.wg.Done()
}

func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitWithTimeout 超时等待退出，返回 true 表示超时，false 表示正常退出
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		w.Wait()
		c <- true
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
