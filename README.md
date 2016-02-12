# Rate Limiting Route Service

## About
A simple in-memory rate limiting route service for Cloud Foundry.

This rate limiter route service app is a forwarding proxy that will limit the number of requests for a given client IP for a given duration.
For example, maximum 10 requests per minute.

If you would like more information about Route Services in Cloud Foundry, please refer to [CF dev docs](http://docs.cloudfoundry.org/services/index.html#route-services).

*NOTE: This is a example only and is not for production use, but enough with the disclaimers...*

## Prerequisites
- A Cloud Foundry and Diego deployment
- CF CLI v6.15+
- an app deployed and running on Cloud Foundry you want to rate limit
- this rate limiter app
```
$ git clone https://github.com/cloudfoundry-samples/ratelimit-service.git
```

- The examples below are using [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
- The app you would like rate limit is running at http://myapp.bosh-lite.com and running on [Diego](https://github.com/cloudfoundry-incubator/diego-release) runtime.

## Install

### Deploy Rate Limiter App
```
$ cd ratelimit-service
$ cf push ratelimiter
```

The rate limiter proxy app will now be running at: https://ratelimiter.bosh-lite.com.


#### Configure limit and duration
To override the default limit (10) and duration in seconds (60), you can set the following application env vars and restage:
```
$ cd set-env ratelimiter rate_duration_in_secs 60
$ cd set-env ratelimiter rate_limit 10
$ cf env ratelimiter

User-Provided:
rate_duration_in_secs: 60
rate_limit: 10

$ cd restage ratelimiter
```

### Create Route Service
The following will create a route service instance using a user-provided service and specifies the route service url (see step above).

```
$ cf create-user-provided-service ratelimiter-service -r https://ratelimiter.bosh-lite.com
Creating user provided service ratelimiter-service in org my-org / space as admin...
OK
```

### Bind Route to Service Instance
The following will create bind the application's route to the route service instance.

```
$ cf bind-route-service bosh-lite.com ratelimiter-service --hostname myapp
Binding route myapp.bosh-lite.com to service instance ratelimiter-service in org my-org / my-space as admin...
OK
```


## Try it out
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



### Unbinding Route Service
If you want to turn off rate limiting, the following will unbind the application's route from the route service instance.

```
cf unbind-route-service bosh-lite.com ratelimiter-service --hostname myapp

Unbinding may leave apps mapped to route myapp.bosh-lite.com vulnerable; e.g. if service instance ratelimiter-service provides authentication. Do you want to proceed?> yes
Unbinding route myapp.bosh-lite.com from service instance ratelimiter-service in org my-org / space my-space as admin...
OK
```

### Logs
You can watch the logs of the rate limiter app to see when requests come it for given IP addresses, as well as when IP address key has expired (i.e. its ok to make more requests)

```
$ cf logs ratelimiter

10:25:00.49-0600 [App/0]  OUT request from [10.244.0.25]
10:25:00.49-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [12/02/2016:16:25:00 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:42192 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:78ebcf18-c546-4d9c-6378-99bf289bb755 response_time:0.004258194 app_id:bfd43944-ae2a-4c65-9ea3-f2aecd390ac1
10:25:00.78-0600 [App/0]  OUT request from [10.244.0.25]
10:25:00.80-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [12/02/2016:16:25:00 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:42198 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:a46616d2-af63-49a3-42ee-f5a8aff4e801 response_time:0.008349286 app_id:bfd43944-ae2a-4c65-9ea3-f2aecd390ac1
10:25:01.00-0600 [App/0]  OUT request from [10.244.0.25]
10:25:01.00-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [12/02/2016:16:25:00 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:42205 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:58a9477f-f59e-4453-61f4-4ef257ecadf4 response_time:0.00661315 app_id:bfd43944-ae2a-4c65-9ea3-f2aecd390ac1
10:25:01.29-0600 [App/0]  OUT request from [10.244.0.25]

10:25:01.29-0600 [App/0]  OUT rate limit exceeded for 10.244.0.25
10:25:01.29-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [12/02/2016:16:25:01 +0000] "HEAD / HTTP/1.1" 429 0 0 "-" "curl/7.43.0" 10.244.0.21:42216 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:b60d0ea1-603f-4a92-6158-e034f3017b41 response_time:0.002778094 app_id:bfd43944-ae2a-4c65-9ea3-f2aecd390ac1
10:25:02.73-0600 [App/0]  OUT request from [10.244.0.25]
10:25:02.73-0600 [App/0]  OUT rate limit exceeded for 10.244.0.25
10:25:02.73-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [12/02/2016:16:25:02 +0000] "HEAD / HTTP/1.1" 429 0 0 "-" "curl/7.43.0" 10.244.0.21:42250 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:33600a1c-5b4a-4067-5617-82cb140c95d6 response_time:0.001626862 app_id:bfd43944-ae2a-4c65-9ea3-f2aecd390ac1
10:25:47.17-0600 [App/0]  OUT removing expired key [10.244.0.25]

```
