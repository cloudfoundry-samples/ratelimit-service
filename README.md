# Rate Limiting Route Service

## About
A simple in-memory rate limiting route service for Cloud Foundry.

This rate limiter route service app is a forwarding proxy that will limit the number of requests/second per client IP address.
For example, maximum 10 requests per second per client.

If you would like more information about Route Services in Cloud Foundry, please refer to [CF dev docs](http://docs.cloudfoundry.org/services/index.html#route-services).

*NOTE: This is an example only and is not for production use, but enough with the disclaimers...*

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


#### (Optional) Configure limit of requests per second
To override the default limit (10), you can set the following application env var and restage:
```
$ cf set-env ratelimiter RATE_LIMIT 1
$ cf env ratelimiter

User-Provided:
RATE_LIMIT: 1

$ cf restage ratelimiter
```

#### (Optional) Skip SSL Validation
If you set the following environment variable to false, the route service
will validate SSL certificates. By default the route service skips SSL validation.

```
$ cf set-env ratelimiter SKIP_SSL_VALIDATION false
$ cf restage ratelimiter
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
To test the rate limiting, you will need to exceed the requests / second limit.
- requests made that are within the limit will return an [HTTP 200](https://httpstatuses.com/200)
- requests made that exceed the limit will return an [HTTP 429 - Too Many Requests](https://httpstatuses.com/429)

**In the examples below we will use a rate limit of 10 requests / per second per client IP**


### Client-side tool
There is a great command line tool (similar to Apache Bench) called [boom](https://github.com/rakyll/boo://github.com/rakyll/boom) which allows you to send a number of requests and also throttle the number of concurrent client requests.

To install client side load testing tool:
```
go get github.com/rakyll/boom
```

### Example (not exceeding rate limit)
For example: 100 requests, 10 concurrently and a QPS of 10

```
$ boom -n 100 -c 10 -q 10 http://myapp.bosh-lite.com
100 / 100 Boooooooooooooooooooooooooooooooooooooooooooooooooooooooooom! 100.00 %

Summary:
  Total:        10.1374 secs.
  Slowest:      0.7535 secs.
  Fastest:      0.0233 secs.
  Average:      0.1478 secs.
  Requests/sec: 9.8644
  Total Data Received:  2500 bytes.
  Response Size per Request:    25 bytes.

Status code distribution:
  [200] 100 responses
```

### Example (exceeding rate limit)
For example: 100 requests, 12 concurrently and a QPS of 12

In this example, since we will be sending 12 requests / sec, then only roughly 84% of the requests should be successful.


```
$ boom -n 100 -c 12 -q 12 http://myapp.bosh-lite.com
100 / 100 Boooooooooooooooooooooooooooooooooooooooooooooooooooooooooom! 100.00 %
Summary:
  Total:        8.3826 secs.
  Slowest:      0.6653 secs.
  Fastest:      0.0352 secs.
  Average:      0.1126 secs.
  Requests/sec: 11.9295
  Total Data Received:  2388 bytes.
  Response Size per Request:    23 bytes.

Status code distribution:
  [200] 86 responses
  [429] 14 responses
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
You can watch the logs of the rate limiter app to see when requests come it for given IP addresses

```
$ cf logs ratelimiter
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.49-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.50-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38101 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:127716be-9cd7-484d-634c-d1051993accb response_time:0.013694649 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.51-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38102 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:cf60dc83-bff1-4e53-6ff7-43a52667e2a4 response_time:0.017450561 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.52-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38109 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:272a34cc-dcaa-4c13-5584-be28e25855a2 response_time:0.027553184 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.52-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38107 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:6c379238-d5fc-4fc5-7b27-e4776d830c41 response_time:0.029788514 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.53-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38106 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:2efdd43d-3429-4869-54f6-634cb4bc9854 response_time:0.038401721 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.53-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38103 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:2f059c95-31b9-42ab-4459-1fe0f54dc9fa response_time:0.038841171 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.54-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38108 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:950c9525-490a-4fd1-7b38-5f74514ea085 response_time:0.046893267 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.55-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38104 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:49c5650b-c511-4fe5-7f56-067e404d7b8b response_time:0.059563104 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.55-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38105 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:de96c743-032f-4977-4466-f9241a03cf6a response_time:0.063431108 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.56-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 200 0 25 "-" "Go-http-client/1.1" 10.244.0.21:38111 x_forwarded_for:"192.168.50.1, 10.244.0.21" x_forwarded_proto:"http" vcap_request_id:081cb579-a038-4511-43e9-b8f2074019b9 response_time:0.070817607 app_id:7a1745bc-d7cb-43a3-8201-c6ac5d75e79c
14:16:20.64-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.64-0600 [App/0]  OUT request from [10.244.0.25]
14:16:20.64-0600 [App/0]  OUT rate limit exceeded for 10.244.0.25
14:16:20.64-0600 [RTR/0]  OUT ratelimiter.bosh-lite.com - [20/02/2016:20:16:20 +0000] "GET / HTTP/1.1" 429 0 17 "-" "Go-http-client/1.1" 10.244.0.21:38144 x_forwarded_for:

```

## Misc
The rate limit app also has a `/stats` endpoint that displays the current list of IPs and available requests, which can be useful for debugging or displaying current stats.

For example, if the limit is set to 10 then the "available" number will be a number between 0 and 10, with zero meaning that the (req/sec) rate limit has been exceeded for the client IP.



```
$ curl ratelimiter.bosh-lite.com/stats
```

```json
[
  {
    "ip": "10.244.0.25",
    "available": 3
  },
  {
    "ip": "10.244.0.29",
    "available": 0
  }
]
```
