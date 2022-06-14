// Package rpc provides an HTTP gateway to a jsonrpc 2.0 service for visibility
// and management of the calling applications L.Loggers.
//
// The HTTP endpoint is a single POST endpoint ending in "/L".  This endpoint
// accepts hmac enveloped jsonrpc 2.0 requests for an rpc service with the
// following methods.
//
// You can find a design doc at https://github.com/scott-cotton/L/blob/main/rpc/design.md.
package rpc
