package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Request struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params []byte `json:"params,omitempty"`
}

func NewRequest[P any](id int, method string, p *P) (*Request, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		return nil, err
	}
	return &Request{
		ID:     id,
		Method: method,
		Params: buf.Bytes(),
	}, nil
}

func Params[P any](req *Request) (*P, error) {
	var params P
	err := json.Unmarshal(req.Params, &params)
	if err != nil {
		return nil, err
	}
	return &params, nil
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []byte `json:"data,omitempty"`
}

type Response struct {
	ID     int    `json:"id"`
	Error  *Error `json:"error,omitempty"`
	Result []byte `json:"result,omitempty"`
}

func NewResponse[Result any](id int, r *Result) (*Response, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}
	return &Response{
		ID:     id,
		Result: buf.Bytes(),
	}, nil
}

func ErrorResponse(id, code int, msg string) *Response {
	fmt.Printf("ErrorResponse %d %d %s\n", id, code, msg)
	return &Response{
		ID: id,
		Error: &Error{
			Code:    code,
			Message: msg,
		},
	}
}

func Result[T any](resp *Response) (*T, error) {
	var res T
	err := json.Unmarshal(resp.Result, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
