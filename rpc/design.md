# L HTTP-RPC

This package uses hmac-sha256 authentication envelop around a jsonrpc-2.0
payload, served under a handler for a POST to a URL ending in "/L".

The service is comprised of 2 methods: 

1. "loggers", a fetch/query method which returns the label mapping for all
   loggers.
1. "apply", a method for applying a configuration using [configuration
   apply](https://pkg.go.dev/github.com/scott-cotton/L#Config.Apply)

## status

This package is under development, some of the functionality below does
not yet exist.

## loggers


Request
```json
{
	"jsonrpc": "2.0",
	"id": <id>,
	"method": "loggers",
	"params": "github.com/.*"
}
```

Response
```json
{
	"jsonrpc": "2.0",
	"id": <id>,
	"result": [
		// index of parent
		"parent": 0, 
		"pkg": "github.com/scott-cotton/L",
		"labels": {
			"a": 10,
			"b": 11
		}
	]
}
```



## apply

Request 
```
{
	"jsonrpc": "2.0",
	"id": <id>,
	"method": "apply",
	"params": {
		"pkgPattern": "github.com/scott-cotton/L": 
		"labels": {
			// add this label only set, do not replace
			"Lset": 1,
			"zebra": 1010
		}
	}
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
The signature is the hmac sha256 signature of the payload (before base64 encoding).

## HTTP gateway

The above is served at an HTTP POST endpoint.  There is no internal gateway,
which unfortunately means that the clients tend to use new connections for
each request.  HTTP gateway was used to make the service accessible via 
standard HTTP tools.




