package rpc

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer("abc", ":4321", "/")
	s.Serve()
}
