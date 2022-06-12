package rpc

import "github.com/scott-cotton/L"

type LoggersRequest struct {
	Request
}

type LoggersResult []L.ConfigNode
