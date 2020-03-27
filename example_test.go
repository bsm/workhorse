package workhorse_test

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/bsm/workhorse"
)

func Example() {
	// define root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count uint32

	// init a worker
	w := workhorse.New(ctx)

	// schedule task "one"
	w.Go("one", func(_ context.Context) error {
		for i := 0; i < 1000; i++ {
			atomic.AddUint32(&count, 1)
		}
		return nil
	})

	// schedule task "two"
	w.Go("two", func(_ context.Context) error {
		for i := 0; i < 1000; i++ {
			atomic.AddUint32(&count, 2)
		}
		return nil
	})

	// wait for both tasks to complete
	if err := w.Wait(); err != nil {
		panic(err)
	}

	fmt.Println(count)
	// Output:
	// 3000
}

func ExampleInstrument() {
	// implement instrumentation
	inst := func(name string, dur time.Duration, err error) {
		if err != nil {
			log.Printf("task %s failed with %v", name, err)
		} else {
			log.Printf("task %s finished in %v", name, dur)
		}
	}

	// a noop task, just waits for 1s
	task := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
		return nil
	}

	// init a worker
	w := workhorse.New(context.Background())

	// run instrumented task every 5s
	w.Go("task", workhorse.Every(
		workhorse.Instrument(task, inst),
		5*time.Second,
	))
	if err := w.Wait(); err != nil {
		panic(err)
	}
}
