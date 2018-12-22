package workerpool

import (
	"context"
	"errors"
)

type WorkerPool struct {
	ctx       context.Context
	scheduler chan func()
}

func New(size int, ctx context.Context) *WorkerPool { //{{{
	pool := &WorkerPool{
		ctx:       ctx,
		scheduler: make(chan func()),
	}
	for i := 0; i < size; i++ {
		go _worker_func(pool.scheduler, pool.ctx)
	}
	return pool
} //}}}

func _worker_func(c chan func(), ctx context.Context) { //{{{
	for {
		select {
		case f := <-c:
			f()
		case <-ctx.Done():
			return
		}
	}
} //}}}
func Exec(p *WorkerPool, f func()) error {
	select {
	case <-p.ctx.Done():
		return errors.New("pool context is cancelled")
	case p.scheduler <- f:
		return nil
	}
}
