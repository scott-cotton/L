package rpc

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	sync.Mutex
	key  []byte
	addr string
	id   int
}

func NewClient(key, addr string) (*Client, error) {
	res := &Client{key: []byte(key), addr: addr}
	return res, nil
}

func (c *Client) Loggers(pat string) (*LoggersResponse, error) {
	c.Lock()
	defer c.Unlock()
	// construct LoggersRequest
	req, err := NewRequest[string](c.getID(), "loggers", &pat)
	if err != nil {
		return nil, fmt.Errorf("error constructing jsonrpc request: %w", err)
	}
	buf := bytes.NewBuffer(nil)
	if err := ToWriter(c.key, buf, req); err != nil {
		panic(err)
	}

	hReq, err := http.NewRequest("POST", c.addr, buf)

	if err != nil {
		return nil, fmt.Errorf("error constructing HTTP request: %w", err)
	}
	hReq.Header.Add("Accept", "application/json")
	hReq.Header.Add("Content-type", "application/json")
	h := &http.Client{}
	resp, err := h.Do(hReq)
	if err != nil {
		return nil, fmt.Errorf("error performing http request: %w", err)
	}
	defer resp.Body.Close()
	gResp, err := FromReader[LoggersResponse](c.key, resp.Body)
	if err != nil {
		if gResp != nil && gResp.Error != nil {
			return nil, errors.New(gResp.Error.Message)
		}
		return nil, fmt.Errorf("error decoding http response: %w", err)
	}
	return gResp, nil
}

func (c *Client) getID() int {
	res := c.id
	c.id++
	return res
}
