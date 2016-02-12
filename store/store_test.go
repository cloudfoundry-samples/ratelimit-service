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

	Describe("Increment", func() {
		BeforeEach(func() {
			store = NewStore()
		})

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
	})

	Describe("ExpiresIn", func() {
		BeforeEach(func() {
			store = NewStore()
		})

		It("expires keys", func() {
			count := store.Increment("foo")
			Expect(count).To(Equal(1))
			Expect(store.CountFor("foo")).To(Equal(1))

			store.ExpiresIn(50*time.Millisecond, "foo")
			time.Sleep(1 * time.Second)
			Expect(store.CountFor("foo")).To(Equal(0))
		})
	})

})
var _ = Describe("Entry", func() {

	Context("new entry", func() {
		It("should not be expired", func() {
			entry := NewEntry()
			Expect(entry.Expired()).To(BeFalse())
		})
	})

	Context("entry with expire time", func() {
		It("should be considered expired", func() {
			entry := NewEntry()
			Expect(entry.Expired()).To(BeFalse())

			entry.Expirable = true
			entry.ExpiryTime = time.Now()
			Eventually(entry.Expired()).Should(BeTrue())
		})
	})
})
