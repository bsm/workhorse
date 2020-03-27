package workhorse_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bsm/workhorse"
)

func Example() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count uint32

	w := workhorse.New(ctx)
	w.Go("one", func(_ context.Context) error {
		for i := 0; i < 1000; i++ {
			atomic.AddUint32(&count, 1)
		}
		return nil
	})
	w.Go("two", func(_ context.Context) error {
		for i := 0; i < 1000; i++ {
			atomic.AddUint32(&count, 2)
		}
		return nil
	})
	if err := w.Wait(); err != nil {
		panic(err)
	}

	fmt.Println(count)

	// Output:
	// 3000
}
