package main

import "context"

type signal interface {
	subscribe(*func(context.Context))
	unsubscribe(*func(context.Context))
}

type Signal[T any] struct {
	value       T
	subscribers set[*func(context.Context)]
}

func (s *Signal[T]) Set(ctx context.Context, value T) {
	s.value = value

	if batch, err := ctx.Value(ctxBatchKey).(set[*func(context.Context)]); err {
		s.subscribers.forEach(func(fn *func(context.Context)) {
			batch.add(fn)
		})
		return
	}

	s.subscribers.forEach(func(fn *func(context.Context)) {
		(*fn)(ctx)
	})
}

func (s *Signal[T]) Get(ctx context.Context) T {
	if deps, ok := ctx.Value(ctxDepsKey).(set[signal]); ok {
		deps.add(s)
		ctx = context.WithValue(ctx, ctxDepsKey, deps)
	}
	return s.value
}

func (s *Signal[T]) subscribe(fn *func(context.Context)) {
	s.subscribers.add(fn)
}

func (s *Signal[T]) unsubscribe(fn *func(context.Context)) {
	s.subscribers.remove(fn)
}

func NewSignal[T any](init T) *Signal[T] {
	return &Signal[T]{value: init, subscribers: newSet[*func(context.Context)]()}
}

const (
	ctxDepsKey  = "signals-deps"
	ctxBatchKey = "signals-batch"
)

func Effect(ctx context.Context, f func(context.Context)) func() {
	deps := newSet[signal]()
	f(context.WithValue(ctx, ctxDepsKey, deps))

	deps.forEach(func(dep signal) {
		dep.subscribe(&f)
	})

	return func() {
		deps.forEach(func(t signal) {
			t.unsubscribe(&f)
		})
	}
}

func Batch(ctx context.Context, f func(context.Context)) {
	var batch set[*func(context.Context)]
	if b, ok := ctx.Value(ctxBatchKey).(set[*func(context.Context)]); ok {
		batch = b
	} else {
		batch = newSet[*func(context.Context)]()
	}

	f(context.WithValue(ctx, ctxBatchKey, batch))

	batch.forEach(func(t *func(context.Context)) {
		(*t)(ctx)
	})
}
