# Rate Limiting Route Service

## About
A simple in-memory rate limiting route service for Cloud Foundry.

This rate limiter route service app is a forwarding proxy that will limit the number of requests for a given client IP for a given duration.
For example, maximum 10 requests per minute.

If you would like more information about Route Services in Cloud Foundry, please refer to [CF dev docs](http://docs.cloudfoundry.org/services/index.html#route-services).

*NOTE: This is a example only and is not for production use, but enough with the disclaimers...*

## Prerequisites
- A Cloud Foundry and Diego deployment
- CF CLI v6.16+
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

10:38:26.42 [App/0] OUT limit [10] duration [1m0s]

10:38:54.86 [App/0] OUT request from [10.244.0.25]
10:38:55.05 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:54 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48163 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:384012cd-dacc-49b5-72dc-0004c65aad56 response_time:0.196000975 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:56.22 [App/0] OUT request from [10.244.0.25]
10:38:56.23 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:56 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48177 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:35516cdb-8e70-4f36-6c97-8eb977df1b1d response_time:0.004578871 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:57.00 [App/0] OUT request from [10.244.0.25]
10:38:57.01 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:57 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48189 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:7dc01bbe-4edf-436f-50a9-7d660a7eeb9d response_time:0.005896061 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:57.44 [App/0] OUT request from [10.244.0.25]
10:38:57.45 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:57 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48199 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:91d9d0fd-108f-4536-48a9-d8f3e9018580 response_time:0.005916887 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:57.91 [App/0] OUT request from [10.244.0.25]
10:38:57.91 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:57 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48206 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:879023c6-7d18-4258-6d8b-a5948e3b6e99 response_time:0.004543364 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:58.36 [App/0] OUT request from [10.244.0.25]
10:38:58.37 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:58 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48213 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:c1089f0d-8af6-483c-7f19-a2ccdb7612ea response_time:0.004303727 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:58.86 [App/0] OUT request from [10.244.0.25]
10:38:58.86 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:58 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48222 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:540b94b8-1c54-4e46-403a-0ad2411a2ece response_time:0.004490439 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:38:59.37 [App/0] OUT request from [10.244.0.25]
10:38:59.37 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:38:59 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48231 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:3b266ca5-329b-4b73-6c1c-12deb3005452 response_time:0.005246971 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:39:00.02 [App/0] OUT request from [10.244.0.25]
10:39:00.03 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:39:00 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48239 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:3979a675-7eeb-46c9-63f5-6f5e69a374ad response_time:0.004491951 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:39:00.70 [App/0] OUT request from [10.244.0.25]
10:39:00.71 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:39:00 +0000] "HEAD / HTTP/1.1" 200 0 0 "-" "curl/7.43.0" 10.244.0.21:48250 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:8d15df2d-571d-4c25-7dd5-024a86802f5e response_time:0.004477301 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:39:02.09 [App/0] OUT request from [10.244.0.25]

10:39:02.09 [App/0] OUT rate limit exceeded for 10.244.0.25
10:39:02.09 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:39:02 +0000] "HEAD / HTTP/1.1" 429 0 0 "-" "curl/7.43.0" 10.244.0.21:48274 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:78b36c9f-eaea-47c9-42f1-cb9595866cc9 response_time:0.001362549 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:39:02.98 [App/0] OUT request from [10.244.0.25]
10:39:02.98 [App/0] OUT rate limit exceeded for 10.244.0.25
10:39:02.98 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:39:02 +0000] "HEAD / HTTP/1.1" 429 0 0 "-" "curl/7.43.0" 10.244.0.21:48284 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:5004bb5b-a716-4a13-61f8-21ae0ec4572b response_time:0.001353071 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55
10:39:05.96 [App/0] OUT request from [10.244.0.25]
10:39:05.96 [App/0] OUT rate limit exceeded for 10.244.0.25
10:39:05.96 [RTR/0] OUT ratelimiter.bosh-lite.com - [12/02/2016:16:39:05 +0000] "HEAD / HTTP/1.1" 429 0 0 "-" "curl/7.43.0" 10.244.0.21:48309 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:bab11bc4-ad22-44ad-4146-151b0b4c04df response_time:0.002815149 app_id:2d0e10f0-3bfc-4fe1-85d7-cc8468cecc55

10:39:54.91 [App/0] OUT removing expired key [10.244.0.25]

10:40:06.41 [App/0] OUT request from [10.244.0.25]


```

## Misc
The rate limit app also has a `/stats` endpoint that displays the current list of IPs and counts, which can be useful for debugging or displaying current stats


```
curl ratelimiter.bosh-lite.com/stats
```

```json
[
  {
    "Ip": "10.244.0.25",
    "Count": 3
  },
  {
    "Ip": "10.244.0.29",
    "Count": 2
  }
]
```
