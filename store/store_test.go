package store_test

import (
	"time"

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

			for i := 1; i < 16; i++ {
				count := store.Increment("bar")
				Expect(count).To(Equal(i))
			}

			Expect(store.CountFor("foo")).To(Equal(10))
			Expect(store.CountFor("bar")).To(Equal(15))

		})

		It("expires keys", func() {
			count := store.Increment("foo")
			Expect(count).To(Equal(1))
			Expect(store.CountFor("foo")).To(Equal(1))
			store.ExpiresIn(50*time.Millisecond, "foo")
			Eventually(store.CountFor("foo")).Should(Equal(0))
		})

	})

})
