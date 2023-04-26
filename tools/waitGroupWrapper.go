package tools

import "sync"

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(fn func()) {
	w.Add(1)
	go func() {
		fn()
		w.Done()
	}()
}
