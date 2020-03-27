package workhorse

import (
	"context"
	"errors"
	"time"
)

// Every applies a task periodically every interval until the first
// failure. Example:
//
// 	w.Run(workhorse.Every(func(ctx context.Context) error {
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
				return ctx.Err()
			case <-t.C:
			}

			if err := task(ctx); err != nil {
				return err
			}
		}
	}
}

// Bypass ignores certain errors.
// Uses errors.Is to make the comparison
//
// 	w.Run(workhorse.Bypass(func(ctx context.Context) error {
// 		select {
// 		case <-ctx.Done()
// 			return ctx.Err() // will be ignored
// 		case <-time.After(time.Hour)
// 			return nil
// 		}
// 	}, context.Canceled, io.EOF))
func Bypass(task Task, ignore ...error) Task {
	return func(ctx context.Context) error {
		err := task(ctx)
		for _, irr := range ignore {
			if errors.Is(err, irr) {
				return nil
			}
		}
		return err
	}
}
