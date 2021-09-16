package timegroup

import (
	"context"
	"errors"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type Group interface {
	Go(f func() error)
	WaitTimeout(time.Duration) error
	Wait() error
}

type group struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	err     error
	errOnce sync.Once
}

func New() Group {
	ctx, cancel := context.WithCancel(context.TODO())
	return &group{ctx: ctx, cancel: cancel}
}

func WithContext(ctx context.Context) (Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &group{ctx: ctx, cancel: cancel}, ctx
}

func (g *group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.recover()
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

func (g *group) recover() {
	if rec := recover(); rec != nil {
		log.Printf("panic:%v\n stack:%v", rec, string(debug.Stack()))
	}
}

func (g *group) WaitTimeout(timeout time.Duration) error {
	go g.Wait()
	select {
	case <-time.After(timeout):
		return errors.New("wait time out")
	case <-g.ctx.Done():
		return g.err
	}
}

func (g *group) Wait() error {
	g.wg.Wait()
	if g.err == nil {
		g.cancel()
	}
	return g.err
}
