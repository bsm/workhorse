# Workhorse

[![Build Status](https://travis-ci.org/bsm/workhorse.svg)](https://travis-ci.org/bsm/workhorse)
[![GoDoc](https://godoc.org/github.com/bsm/workhorse?status.png)](http://godoc.org/github.com/bsm/workhorse)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A simple worker abstraction on top of [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) with custom middlewares.

## Examples

```go
import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bsm/workhorse"
)

func main() {
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
	// 3000}
```

## Documentation

Full documentation is available on [GoDoc](https://pkg.go.dev/github.com/bsm/workhorse)
