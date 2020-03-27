package workhorse

import (
	"context"
	"time"
)

// Every applies a task periodically every interval until the first
// failure. Example:
//
// 	w.Go("task", workhorse.Every(func(ctx context.Context) error {
// 		fmt.Println("still alive!")
// 		return nil
// 	}, time.Minute))
func Every(task Task, interval time.Duration) Task {
	return func(ctx context.Context) error {
		t := time.NewTicker(interval)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-t.C:
			}

			if err := task(ctx); err != nil {
				return err
			}
		}
	}
}

// Retry retries a task on failures with a linear backoff.
// Set the numRetries to -1 to retry forever. Example:
//
// 	w.Go("task", workhorse.Retry(func(ctx context.Context) error {
// 		return fmt.Errorf("instant failure!")
// 	}, 4, time.Second))
func Retry(task Task, numRetries int, backoff time.Duration) Task {
	return func(ctx context.Context) error {
		t := time.NewTimer(backoff)
		defer t.Stop()

		remaining := numRetries
		for {
			err := task(ctx)
			if err == nil {
				return nil
			} else if remaining == 0 {
				return err
			} else if remaining > 0 {
				remaining--
			}

			t.Reset(backoff)

			select {
			case <-ctx.Done():
				return nil
			case <-t.C:
			}
		}
	}
}

// Instrument allows to instrument tasks.
func Instrument(task Task, tfn func(name string, runTime time.Duration, err error)) Task {
	return func(ctx context.Context) error {
		start := time.Now()
		err := task(ctx)
		elapsed := time.Since(start)
		tfn(TaskName(ctx), elapsed, err)
		return err
	}
}
