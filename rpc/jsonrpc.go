package rpc

import (
	"bytes"
	"encoding/json"
)

type Request[P any] struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params []byte `json:"params,omitempty"`
}

func NewRequest[P any](id int, method string, p *P) (*Request[P], error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		return nil, err
	}
	return &Request[P]{
		ID:     id,
		Method: method,
		Params: buf.Bytes(),
	}, nil
}

func (r *Request[P]) GetParams() (*P, error) {
	var params P
	err := json.Unmarshal(r.Params, &params)
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

type Response[Result any] struct {
	ID     int    `json:"id"`
	Error  *Error `json:"error,omitempty"`
	Result []byte `json:"result,omitempty"`
}

func NewResponse[Result any](id int, method string, r *Result) (*Response[Result], error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}
	return &Response[Result]{
		ID:     id,
		Result: buf.Bytes(),
	}, nil
}

func ErrorResponse(id, code int, msg string) *Response[any] {
	return &Response[any]{
		ID: id,
		Error: &Error{
			Code:    code,
			Message: msg,
		},
	}
}

func (r *Response[Result]) GetResult() (*Result, error) {
	var res Result
	err := json.Unmarshal(r.Result, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
