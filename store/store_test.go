package store_test

import (
	. "github.com/cloudfoundry-samples/ratelimit-service/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {

	var (
		store Store
	)

	BeforeEach(func() {
		store = NewStore()
	})

	Context("non-concurrently", func() {

		It("increments", func() {

			for i := 1; i < 11; i++ {
				count := store.Increment("foo")
				Expect(count).To(Equal(i))
			}

			for i := 1; i < 11; i++ {
				count := store.Increment("bar")
				Expect(count).To(Equal(i))
			}

		})

	})

})
