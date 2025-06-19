package utils

import (
	"sync"
	"sync/atomic"
)

type OnceV2 struct {
	m    sync.Mutex
	done uint32
}

func (o *OnceV2) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 { //fast path
		return nil
	}
	return o.slowDo(f)
}

func (o *OnceV2) slowDo(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = f()
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}
