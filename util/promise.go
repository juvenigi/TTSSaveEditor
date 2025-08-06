package util

import (
	"context"
	"fmt"
	"sync"
)

type Result[T any] struct {
	Result T
	Err    error
}

type Promise[T any] struct {
	once sync.Once
	done chan struct{}
	val  *T
	err  error
}

func (p *Promise[T]) deadGopherSwitch() {
	err := recover()
	if err == nil {
		return
	}
	switch v := err.(type) {
	case error:
		p.reject(v)
	default:
		p.reject(fmt.Errorf("%+v", v))
	}
}

func RunAsync[T any](ctx context.Context, fn func(ctx context.Context) (T, error)) *Promise[T] {
	p := &Promise[T]{done: make(chan struct{})}
	go func() {
		defer p.deadGopherSwitch()
		val, err := fn(ctx)

		// defence: only one goroutine will resolve or reject the promise
		p.once.Do(func() {
			if err == nil {
				p.val = &val
			} else {
				p.err = err
			}
			close(p.done)
		})
	}()
	return p
}

func (p *Promise[T]) Await(ctx context.Context) (*T, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
		return p.val, p.err
	}
}

// NewResolvers returns a promise and its resolver/rejector
func NewResolvers[T any]() (*Promise[T], func(T), func(error)) {
	p := &Promise[T]{done: make(chan struct{})}
	return p, p.resolve, p.reject
}

func (p *Promise[T]) resolve(v T) {
	p.once.Do(func() {
		p.val = &v
		close(p.done)
	})
}

func (p *Promise[T]) reject(err error) {
	p.once.Do(func() {
		p.err = err
		close(p.done)
	})
}
