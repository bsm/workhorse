package workhorse

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Task is a task function.
type Task func(context.Context) error

// Worker runs jobs.
type Worker struct {
	ctx  context.Context
	errs *errgroup.Group
}

// New returns a Worker under a global ctx.
func New(ctx context.Context) *Worker {
	errs, ectx := errgroup.WithContext(ctx)
	return &Worker{ctx: ectx, errs: errs}
}

// Go starts a named background task.
func (w *Worker) Go(name string, task func(ctx context.Context) error) {
	ctx := context.WithValue(w.ctx, taskNameCtxKey, name)
	w.errs.Go(func() error {
		return task(ctx)
	})
}

// Wait blocks and waits for jobs to complete and returns the first error (if any).
func (w *Worker) Wait() error {
	return w.errs.Wait()
}
