package workhorse_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/bsm/workhorse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	var subject *workhorse.Worker
	var bg = context.Background()

	BeforeEach(func() {
		subject = workhorse.New(bg)
	})

	It("should run jobs", func() {
		var n uint32
		subject.Go("one", func(_ context.Context) error {
			for i := 0; i < 1000; i++ {
				atomic.AddUint32(&n, 1)
			}
			return nil
		})

		subject.Go("two", func(_ context.Context) error {
			for i := 0; i < 1000; i++ {
				atomic.AddUint32(&n, 2)
			}
			return nil
		})

		Expect(subject.Wait()).To(Succeed())
		Expect(n).To(Equal(uint32(3000)))

		Expect(subject.Wait()).To(Succeed())
		Expect(n).To(Equal(uint32(3000)))

		subject.Go("three", func(_ context.Context) error {
			for i := 0; i < 111; i++ {
				atomic.AddUint32(&n, 3)
			}
			return nil
		})
		Expect(subject.Wait()).To(Succeed())
		Expect(n).To(Equal(uint32(3333)))
	})

	It("should extract task names from context", func() {
		subject.Go("one", func(ctx context.Context) error {
			Expect(workhorse.TaskName(ctx)).To(Equal("one"))
			return nil
		})
		subject.Go("two", func(ctx context.Context) error {
			Expect(workhorse.TaskName(ctx)).To(Equal("two"))
			return nil
		})
		Expect(subject.Wait()).To(Succeed())
		Expect(workhorse.TaskName(bg)).To(Equal(""))
	})

	It("should handle failures and wait for all tasks to finish", func() {
		var n uint32
		subject.Go("one", func(ctx context.Context) error {
			<-ctx.Done()
			for i := 0; i < 1000; i++ {
				atomic.AddUint32(&n, 1)
			}
			return nil
		})
		subject.Go("two", func(ctx context.Context) error {
			return errors.New("two failed")
		})
		Expect(subject.Wait()).To(MatchError("two failed"))
		Expect(n).To(Equal(uint32(1000)))
	})

	It("should not raise parent context failures", func() {
		parent, cancel := context.WithCancel(bg)
		w := workhorse.New(parent)
		w.Go("one", func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})
		cancel()
		Expect(w.Wait()).To(Succeed())
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "workhorse")
}
