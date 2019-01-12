package workerpool

import (
	"context"
	"errors"
	"sync"
)

type WorkerPool struct {
	ctx       context.Context
	scheduler chan func()
}

func New(size int, ctx context.Context) *WorkerPool { //{{{
	var wg sync.WaitGroup
	pool := &WorkerPool{
		ctx:       context.WithValue(ctx, "wg", &wg),
		scheduler: make(chan func()),
	}
	for i := 0; i < size; i++ {
		go _worker_func(pool.scheduler, pool.ctx)
	}
	return pool
} //}}}

func _worker_func(c chan func(), ctx context.Context) { //{{{
	wg, ok := ctx.Value("wg").(*sync.WaitGroup)
	for {
		select {
		case f := <-c:
			f()
			if ok {
				wg.Done()
			}

		case <-ctx.Done():
			return
		}
	}
} //}}}
func Exec(p *WorkerPool, f func()) error {
	wg, ok := p.ctx.Value("wg").(*sync.WaitGroup)
	if ok {
		wg_add(p)
	}
	select {
	case <-p.ctx.Done():
		if ok {
			wg.Done()
		}
		return errors.New("pool context is cancelled")
	case p.scheduler <- f:
		return nil
	}
}
func Wait(p *WorkerPool) { //{{{
	wg, ok := p.ctx.Value("wg").(*sync.WaitGroup)
	if ok {
		wg.Wait()
	}
} //}}}
func wg_add(p *WorkerPool) {
	wg, ok := p.ctx.Value("wg").(*sync.WaitGroup)
	if ok {
		wg.Add(1)
	}
}
