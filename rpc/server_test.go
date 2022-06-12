package rpc

import (
	"os"
	"testing"
)

// interactive use with client_test.go
func TestServer(t *testing.T) {
	if os.Getenv("LSERVE") == "" {
		t.Logf("set LSERVE=1 to run")
		return
	}
	s := NewServer("abc", ":4321", "/")
	s.Serve()
}
