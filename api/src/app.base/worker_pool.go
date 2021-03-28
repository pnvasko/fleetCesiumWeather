package appcontroler

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants"
)

type WorkerPool struct {
	pool *ants.Pool
	log  *Logger
}

func NewWorkerPool(ctx context.Context, poolSize int) (*WorkerPool, error) {
	var err error
	wp := &WorkerPool{}

	wp.log = GetLog(ctx, wp)

	if wp.pool, err = ants.NewPool(poolSize, ants.WithPreAlloc(true)); err != nil {
		errmsg := fmt.Errorf("NewWorkerPool error init WorkerPool: %s", err)
		return nil, errmsg
	}

	return wp, nil
}

func (wp *WorkerPool) Release() {
	wp.pool.Release()
}

func (wp *WorkerPool) Submit(f func()) {
	if err := wp.pool.Submit(f); err != nil {
		wp.log.Error("Error submit: ", err)
	}
	return
}
