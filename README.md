# Rate Limiting Route Service

## About
A simple in-memory rate limiting route service for Cloud Foundry.

This rate limiter app is a forwarding proxy that will rate limit the number of requests for a given client IP for a given duration
(i.e. max 10 requests per minute).

If you would like more information about Route Services in Cloud Foundry, please refer to this [doc](https://docs.google.com/document/d/1bGOQxiKkmaw6uaRWGd-sXpxL0Y28d3QihcluI15FiIA/edit#heading=h.8djffzes9pnb).

*NOTE: This is a example only and is not for production use, but enough with the disclaimers...*

## Prerequisites
- A Cloud Foundry
- an application deployed and running on Cloud Foundry
- CF CLI v6.15+

Note: The examples below are assuming using bosh-lite and assuming the application you would like rate limit is http://myapp.bosh-lite.com.

## Install

### Deploy Rate Limiter App
    $ git clone https://github.com/cloudfoundry-samples/ratelimit-service.git
    $ cd ratelimit-service
    $ cf push ratelimiter

The rate limiter proxy app will now be running at: https://ratelimiter.bosh-lite.com.

### Create Route Service
The following will create a route service instance using a user-provided service and specifies the route service url (see step above).

    $ cf create-user-provided-service ratelimit-service -r https://ratelimiter.bosh-lite.com

### Bind Route to Service Instance
The following will create bind the application's route to the route service instance.

    $ cf bind-route-service bosh-lite.com ratelimiter-service --hostname myapp


## Testing
To test the rate limiting, make more than 10 requests to the application within a minute. The first 10 requests will return an HTTP 200 and any subsequent requests over 10 within 60 seconds will return HTTP 429 - Too Many Requests.

    $ curl http://myapp.bosh-lite.com
    
