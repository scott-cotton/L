package rpc

import (
	"encoding/json"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewClient("abc", "http://localhost:4321/L")
	if err != nil {
		t.Error(err)
		return
	}
	lr, err := client.Loggers(".*")
	if err != nil {
		t.Error(err)
		return
	}
	res, err := lr.GetResult()
	if err != nil {
		t.Logf("%#v", lr)
		t.Error(err)
		return
	}
	jenc := json.NewEncoder(os.Stdout)
	jenc.SetIndent("", "  ")
	jenc.Encode(res)
}
