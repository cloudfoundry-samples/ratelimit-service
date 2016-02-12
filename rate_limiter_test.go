package main_test

import (
	. "github.com/cloudfoundry-samples/ratelimit-service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RateLimiter", func() {

	var (
		limiter *RateLimiter
	)

	BeforeEach(func() {
		limiter = NewRateLimiter(10)
	})

	It("rate limits", func() {

		Expect(1).To(Equal(1))
	})

})
