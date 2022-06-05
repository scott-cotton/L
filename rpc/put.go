package rpc

type PutRequest struct {
	Labels map[string]bool
}

type PutResponse struct {
	Added   map[string]int
	Removed map[string]int
}
