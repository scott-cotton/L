# L HTTP-RPC

This package uses hmac-sha256 authentication envelop around a jsonrpc-2.0
payload, served under a handler for a POST to a URL ending in "/L".

The service is comprised of 2 methods: 

1. "loggers", a fetch/query method which returns the label mapping for all
   loggers.
1. "apply", a method for applying a configuration using [configuration
   apply](https://pkg.go.dev/github.com/scott-cotton/L#Config.Apply)


## loggers


Request
```json
{
	"jsonrpc": "2.0",
	"id": 123,
	"method": "loggers",
}
```

Response
```json
{
	"jsonrpc": "2.0",
	"id": 123,
	"result": [
		{
			"parent": 0, 
			"pkg": "github.com/scott-cotton/L",
			"labels": {
				"a": 10,
				"b": 11
			}
		}
	]
}
```

The parent field in each array entry of the result is the index in the array
of the parent of the logger, -1 if there is no parent (the root).



## apply

Request 
```
{
	"jsonrpc": "2.0",
	"id": 456,
	"method": "apply",
	"params": {
		"pkgPattern": "github.com/scott-cotton/L": 
		"opts": {
			"removeAbsentLabels": true
		},
		"config": {
			"labels": {
				"zebra": 1010
			}
		}
	}
}
```

- pkgPattern indicates which packages to match.
- opts is an "github.com/scott-cotton/L".ApplyOpts object, but it is always recursive.
- config is a configuration object.  Currrently, this contains only labels.  Later
we will consider adding formatters, writers, and middleware via registration.

Response

The response is in the same form as a loggers response, except each result
item does not have a "parent" field and the array does not represent a tree.
Instead, the array contains each modified configuration.
```json
{
	"jsonrpc": "2.0",
	"id": 456,
	"result": [
		{
			"pkg": "github.com/scott-cotton/L",
			"labels": {
				"a": 10,
				"b": 11
			}
		}
	]
}
```

## hmac envelop

Given the requests and responses above, we wrap them in signed payloads where
the client and server are assumed to have shared a secret key out of band.

Each message takes the form

```json
{
	"payload": "xkd003=",
	"signature": "kdoj="
}
```

The payload is the base64 encoded byte array of the respective messages above.
The signature is the base64 encoding of the hmac sha256 signature of the payload 
(before base64 encoding).

## HTTP gateway

The above is served at an HTTP POST endpoint.  There is no internal gateway,
which unfortunately means that the clients tend to use new connections for
each request.  HTTP gateway was used to make the service accessible via 
standard HTTP tools.




