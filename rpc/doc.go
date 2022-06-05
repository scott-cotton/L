// Package rpc provides an HTTP gateway to a jsonrpc 2.0 service
// for visibility and management of the calling applications L.Loggers.
//
// The HTTP endpoint is a single POST endpoint ending in "/L".  This
// endpoint accepts jsonrpc 2.0 requests for an rpc service with the
// following methods.
//
//  - "loggers"  The loggers method retrieves information about labels
//  for every active logger at the point in time in which the call is
//  processed.
//
package rpc
