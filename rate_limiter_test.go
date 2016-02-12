package main_test

import (
	"time"

	. "github.com/cloudfoundry-samples/ratelimit-service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RateLimiter", func() {

	var (
		limit    int
		duration time.Duration
		limiter  *RateLimiter
	)

	Describe("ExceedsLimit", func() {

		BeforeEach(func() {
			limit = 10
			duration = 1 * time.Second
			// only max of n per duration
			limiter = NewRateLimiter(limit, duration)
		})

		It("reports if rate exceeded and resets", func() {
			for i := 0; i < limit; i++ {
				Expect(limiter.ExceedsLimit("192.168.1.1")).To(BeFalse())
			}
			Expect(limiter.ExceedsLimit("192.168.1.1")).To(BeTrue())

			// wait duration to ensure reset
			time.Sleep(duration + 500*time.Millisecond)
			Expect(limiter.ExceedsLimit("192.168.1.1")).To(BeFalse())
		})
	})

})
