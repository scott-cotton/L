package rpc

type LoggersRequest struct {
	Request[string]
}

type LoggersResponse struct {
	Response[map[string]map[string]int]
}
