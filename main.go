package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	DEFAULT_PORT            = "8080"
	CF_FORWARDED_URL_HEADER = "X-Cf-Forwarded-Url"
)

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
			forwardedURL := req.Header.Get(CF_FORWARDED_URL_HEADER)

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
	transport http.RoundTripper
}

func newRateLimiter() *RateLimiter {
	return &RateLimiter{
		transport: http.DefaultTransport,
	}
}

func (r *RateLimiter) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response

	// TODO: add simple in-memory rate limiting logic

	res, err = r.transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	return res, err
}
