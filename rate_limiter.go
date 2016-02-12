package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-samples/ratelimit-service/store"
)

type RateLimiter struct {
	limit int
	store store.Store
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		limit: limit,
		store: store.NewStore(),
	}
}

func (r *RateLimiter) ExceedsLimit(ip string) bool {
	current := r.store.Increment(ip)

	// if first request set expiry time
	if current == 1 {
		r.store.ExpiresIn(60*time.Second, ip)
	}

	// if exceeds limit
	if current > limit {
		fmt.Printf("rate limit exceeded for %s\n", ip)
		return true
	}

	return false
}
