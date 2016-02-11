package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	DEFAULT_PORT     = "8080"
	CF_FORWARDED_URL = "X-Cf-Forwarded-Url"
	REMOTE_ADDRESS   = "REMOTE_ADDR"

	// TODO: make configurable
	limit = 10
)

type Store interface {
	Increment(string) int
	ExpiresIn(int, string)
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	log.SetOutput(os.Stdout)

	log.Fatal(http.ListenAndServe(":"+port, newProxy()))
}

func newProxy() http.Handler {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(CF_FORWARDED_URL)

			url, err := url.Parse(forwardedURL)
			if err != nil {
				log.Fatalln(err.Error())
			}

			req.URL = url
			req.Host = url.Host
		},
		Transport: newRateLimiter(),
	}
	return proxy
}

type RateLimiter struct {
	store     Store
	transport http.RoundTripper
}

func newRateLimiter() *RateLimiter {
	return &RateLimiter{
		store:     &InMemoryStore{},
		transport: http.DefaultTransport,
	}
}

func (r *RateLimiter) exceedsLimit(ip string) bool {
	// need to implement rate limiting logic

	current := r.store.Increment(ip)

	// if exceeds limit
	if current > limit {
		fmt.Printf("rate limit exceeded for %s\n", ip)
	}

	// if first request set expiry time
	if current == 1 {
		r.store.ExpiresIn(60, ip)
	}

	return false
}

func (r *RateLimiter) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response

	// TODO: add simple in-memory rate limiting logic
	remoteIP := req.Header.Get(REMOTE_ADDRESS)
	if r.exceedsLimit(remoteIP) {
		// fix this to properly return an http status of 429
		return nil, errors.New("http 429 - too many requests")
	}

	res, err = r.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, err
}

type InMemoryStore struct {
}

func (s *InMemoryStore) Increment(key string) int {
	return 1
}
func (s *InMemoryStore) ExpiresIn(secs int, key string) {
}
