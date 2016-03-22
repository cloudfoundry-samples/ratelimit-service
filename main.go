package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	DEFAULT_PORT     = "8080"
	CF_FORWARDED_URL = "X-Cf-Forwarded-Url"
	DEFAULT_LIMIT    = 10
)

var (
	limit       int
	rateLimiter *RateLimiter
)

func main() {
	log.SetOutput(os.Stdout)

	limit = getEnv("RATE_LIMIT", DEFAULT_LIMIT)
	log.Printf("limit per sec [%d]\n", limit)

	rateLimiter = NewRateLimiter(limit)

	http.HandleFunc("/stats", statsHandler)
	http.Handle("/", newProxy())
	log.Fatal(http.ListenAndServe(":"+getPort(), nil))
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
		Transport: newRateLimitedRoundTripper(),
	}
	return proxy
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := json.Marshal(rateLimiter.GetStats())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, string(stats))
}

func getPort() string {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	return port
}

func skipSslValidation() bool {
	var skipSslValidation bool
	var err error
	if skipSslValidation, err = strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION")); err != nil {
		skipSslValidation = true
	}
	return skipSslValidation
}
func getEnv(env string, defaultValue int) int {
	var (
		v      string
		config int
	)
	if v = os.Getenv(env); len(v) == 0 {
		return defaultValue
	}

	config, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return config
}

type RateLimitedRoundTripper struct {
	rateLimiter *RateLimiter
	transport   http.RoundTripper
}

func newRateLimitedRoundTripper() *RateLimitedRoundTripper {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSslValidation()},
	}
	return &RateLimitedRoundTripper{
		rateLimiter: rateLimiter,
		transport:   tr,
	}
}

func (r *RateLimitedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response

	remoteIP := strings.Split(req.RemoteAddr, ":")[0]

	log.Printf("request from [%s]\n", remoteIP)
	if r.rateLimiter.ExceedsLimit(remoteIP) {
		resp := &http.Response{
			StatusCode: 429,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Too many requests")),
		}
		return resp, nil
	}

	res, err = r.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
