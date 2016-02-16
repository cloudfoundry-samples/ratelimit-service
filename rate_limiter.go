package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-samples/ratelimit-service/store"
)

type Stats []Stat
type Stat struct {
	Ip    string
	Count int
}

type RateLimiter struct {
	limit    int
	duration time.Duration
	store    store.Store
}

func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		duration: duration,
		store:    store.NewStore(),
	}
}

func (r *RateLimiter) ExceedsLimit(ip string) bool {
	current := r.store.Increment(ip)

	// if first request set expiry time
	if current == 1 {
		r.store.ExpiresIn(r.duration, ip)
	}

	// if exceeds limit
	if current > r.limit {
		fmt.Printf("rate limit exceeded for %s\n", ip)
		return true
	}

	return false
}

func (r *RateLimiter) GetStats() Stats {
	s := Stats{}
	for k, v := range r.store.Current() {
		s = append(s, Stat{
			Ip:    k,
			Count: v,
		})
	}
	return s
}
