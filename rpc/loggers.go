package rpc

import "github.com/scott-cotton/L"

type LoggersRequest struct {
	Request[struct{}]
}

type LoggersResponse struct {
	Response[[]L.ConfigNode]
}
