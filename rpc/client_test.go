package rpc

import (
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewClient("abc", "http://localhost:4321/L")
	if err != nil {
		t.Error(err)
		return
	}
	lr, err := client.Loggers()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := lr.GetResult()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v\n", res)
}
