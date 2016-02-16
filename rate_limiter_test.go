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
			ip := "192.168.1.1"

			for i := 0; i < limit; i++ {
				Expect(limiter.ExceedsLimit(ip)).To(BeFalse())
			}
			Expect(limiter.ExceedsLimit(ip)).To(BeTrue())

			// wait duration to ensure reset
			time.Sleep(duration + 500*time.Millisecond)
			Expect(limiter.ExceedsLimit(ip)).To(BeFalse())
		})
	})

	Describe("Stats", func() {
		BeforeEach(func() {
			limit = 10
			duration = 5 * time.Second
			limiter = NewRateLimiter(limit, duration)
		})

		It("reports stats ", func() {

			for i := 5; i < limit; i++ {
				ip := "192.168.1.100"
				Expect(limiter.ExceedsLimit(ip)).To(BeFalse())
			}
			for i := 7; i < limit; i++ {
				ip := "192.168.1.101"
				Expect(limiter.ExceedsLimit(ip)).To(BeFalse())
			}

			stats := limiter.GetStats()
			Expect(len(stats)).To(Equal(2))
		})

	})

})
