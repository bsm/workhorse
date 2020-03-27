package workhorse_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/bsm/workhorse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Every", func() {
	var bg = context.Background()

	It("should run until the first error", func() {
		var n uint32

		w := workhorse.New(bg)
		w.Go("repeat", workhorse.Every(func(_ context.Context) error {
			atomic.AddUint32(&n, 1)
			return fmt.Errorf("failure")
		}, time.Microsecond))
		Expect(w.Wait()).To(MatchError("failure"))
		Expect(n).To(Equal(uint32(1)))
	})

	It("should stop when context is canceled", func() {
		var n uint32

		ctx, cancel := context.WithCancel(bg)
		defer cancel()

		w := workhorse.New(ctx)
		w.Go("repeat", workhorse.Every(func(_ context.Context) error {
			if atomic.AddUint32(&n, 1) == 10 {
				cancel()
			}
			return nil
		}, time.Microsecond))
		Expect(w.Wait()).To(Succeed())
		Expect(n).To(Equal(uint32(10)))
	})
})

var _ = Describe("Retry", func() {
	var bg = context.Background()

	It("should retry", func() {
		var n uint32

		w := workhorse.New(bg)
		w.Go("repeat", workhorse.Retry(func(_ context.Context) error {
			atomic.AddUint32(&n, 1)
			return fmt.Errorf("failure")
		}, 3, time.Microsecond))
		Expect(w.Wait()).To(MatchError("failure"))
		Expect(n).To(Equal(uint32(4)))
	})

	It("should not retry if 0", func() {
		var n uint32

		w := workhorse.New(bg)
		w.Go("repeat", workhorse.Retry(func(_ context.Context) error {
			atomic.AddUint32(&n, 1)
			return fmt.Errorf("failure")
		}, 0, time.Microsecond))
		Expect(w.Wait()).To(MatchError("failure"))
		Expect(n).To(Equal(uint32(1)))
	})

	It("may retry forever", func() {
		var n uint32

		ctx, cancel := context.WithCancel(bg)
		defer cancel()

		w := workhorse.New(ctx)
		w.Go("repeat", workhorse.Retry(func(_ context.Context) error {
			if atomic.AddUint32(&n, 1) == 10 {
				cancel()
			}
			return fmt.Errorf("failure")
		}, -1, time.Microsecond))
		Expect(w.Wait()).To(Succeed())
		Expect(n).To(Equal(uint32(10)))
	})
})
