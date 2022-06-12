package rpc

import (
	"encoding/json"
	"github.com/scott-cotton/L"
	"os"
	"testing"
)

func TestClientApply(t *testing.T) {
	if os.Getenv("LSERVE") == "" {
		t.Logf("set LSERVE=1 to run")
		return
	}
	client, err := NewClient("abc", "http://localhost:4321/L")
	if err != nil {
		t.Error(err)
		return
	}
	parms := &ApplyParams{
		PkgPattern: ".*",
		Config: &L.Config{
			Labels: map[string]int{
				"hello": 3,
			},
		},
	}
	ar, err := client.Apply(parms)
	if err != nil {
		t.Error(err)
		return
	}
	jenc := json.NewEncoder(os.Stdout)
	jenc.SetIndent("", "  ")
	jenc.Encode(ar)
}

func TestClientLogger(t *testing.T) {
	if os.Getenv("LSERVE") == "" {
		t.Logf("set LSERVE=1 to run")
		return
	}
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
	jenc := json.NewEncoder(os.Stdout)
	jenc.SetIndent("", "  ")
	jenc.Encode(lr)
}
