# Rate Limiting Route Service

## About
A simple in-memory rate limiting route service for Cloud Foundry.

This rate limiter route service app is a forwarding proxy that will limit the number of requests for a given client IP for a given duration.
For example, maximum 10 requests per minute.

If you would like more information about Route Services in Cloud Foundry, please refer to this [doc](https://docs.google.com/document/d/1bGOQxiKkmaw6uaRWGd-sXpxL0Y28d3QihcluI15FiIA/edit#heading=h.8djffzes9pnb).

*NOTE: This is a example only and is not for production use, but enough with the disclaimers...*

## Prerequisites
- A Cloud Foundry and Diego deployment
- CF CLI v6.15+
- an app deployed and running on Cloud Foundry you want to rate limit
- this rate limiter app
      $ git clone https://github.com/cloudfoundry-samples/ratelimit-service.git

Notes:
- The examples below are using [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
- The app you would like rate limit is running at http://myapp.bosh-lite.com and running on [Diego](https://github.com/cloudfoundry-incubator/diego-release) runtime.

## Install

### Deploy Rate Limiter App
```
$ cd ratelimit-service
$ cf push ratelimiter
```

The rate limiter proxy app will now be running at: https://ratelimiter.bosh-lite.com.

### Create Route Service
The following will create a route service instance using a user-provided service and specifies the route service url (see step above).

```
$ cf create-user-provided-service ratelimiter-service -r https://ratelimiter.bosh-lite.com
Creating user provided service ratelimiter-service in org o / space s as admin...
OK
```

### Bind Route to Service Instance
The following will create bind the application's route to the route service instance.

```
$ cf bind-route-service bosh-lite.com ratelimiter-service --hostname myapp
Binding route myapp.bosh-lite.com to service instance ratelimiter-service in org my-org / my-space s as admin...
OK
```


## Testing
To test the rate limiting, make more than 10 requests to the application within a minute. The first 10 requests will return an [HTTP 200](https://httpstatuses.com/200) and any subsequent requests over 10 within 60 seconds will return [HTTP 429 - Too Many Requests](https://httpstatuses.com/429).

```
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 200 OK
$ curl -I myapp.bosh-lite.com
HTTP/1.1 429 Too Many Requests
$ curl -I myapp.bosh-lite.com
HTTP/1.1 429 Too Many Requests
. . .
```
