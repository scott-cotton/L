package rpc

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/scott-cotton/L"
)

func (s *Server) Handler(resp http.ResponseWriter, req *http.Request) {
	s.log.Dict().
		Field("remote", req.RemoteAddr).
		Field("method", req.Method).
		Log()
	switch req.Method {
	case "POST":
		s.ServiceHandler(resp, req)
	default:
		s.warn.Dict().
			Field("bad request method", req.Method).
			Log()
		http.Error(resp, "invalid request method", http.StatusBadRequest)
	}
}

func (s *Server) ServiceHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	r, err := fromReader[Request](s.key, req.Body)
	if err != nil {
		if r == nil {
			s.HTTPError(w, err)
			return
		}
		resp := ErrorResponse(r.ID, 1, err.Error())
		if err := toWriter(s.key, w, resp); err != nil {
			s.HTTPError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	switch r.Method {
	case "loggers":
		tree := LoggersResult(L.ConfigTree())
		rpcResp, err := NewResponse[LoggersResult](r.ID, &tree)
		if err != nil {
			s.JSONRPCError(w, r.ID, 3, err)
			return
		}
		if err := toWriter(s.key, w, rpcResp); err != nil {
			resp := ErrorResponse(r.ID, 2, err.Error())
			if err := toWriter(s.key, w, resp); err != nil {
				s.HTTPError(w, err)
			}
			return
		}
		return
	case "apply":
		applyParams, err := Params[ApplyParams](r)
		if err != nil {
			s.JSONRPCError(w, r.ID, 3, err)
			return
		}
		result, err := Apply(applyParams)
		if err != nil {
			s.JSONRPCError(w, r.ID, 3, err)
			return
		}
		resp, err := NewResponse[ApplyResult](r.ID, &result)
		if err != nil {
			s.JSONRPCError(w, r.ID, 3, err)
			return
		}
		if err := toWriter(s.key, w, resp); err != nil {
			s.HTTPError(w, err)
		}

	default:
		// -32601 is from jsonrpc 2.0
		resp := ErrorResponse(r.ID, -32601,
			fmt.Sprintf("%q: invalid method", r.Method))
		if err := toWriter(s.key, w, resp); err != nil {
			s.HTTPError(w, err)
		}
	}
}

func (s *Server) JSONRPCError(w http.ResponseWriter, id, code int, err error) {
	w.WriteHeader(http.StatusOK)
	resp := ErrorResponse(id, code, err.Error())
	if err := toWriter(s.key, w, resp); err != nil {
		s.HTTPError(w, err)
	}
}

func (s *Server) HTTPError(w http.ResponseWriter, err error) {
	s.warn.Dict().Err(err).Log()
	status := http.StatusInternalServerError
	if errors.Is(err, ErrSignature) {
		status = http.StatusUnauthorized
	}
	http.Error(w, err.Error(), status)
}
