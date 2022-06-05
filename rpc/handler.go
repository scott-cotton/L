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
	r, err := FromReader[Request[any]](s.key, req.Body)
	if err != nil {
		if r == nil {
			s.HTTPError(w, err)
			return
		}
		resp := ErrorResponse(r.ID, 1, err.Error())
		if err := ToWriter(s.key, w, resp); err != nil {
			s.HTTPError(w, err)
		}
		return
	}
	switch r.Method {
	case "loggers":
		lr := (*Request[[]string])(r)
		_ = lr
		parms, err := lr.GetParams()
		if err != nil {
			w.WriteHeader(http.StatusOK)
			resp := ErrorResponse(r.ID, 1, err.Error())
			if err := ToWriter(s.key, w, resp); err != nil {
				s.HTTPError(w, err)
			}
			return
		}
		res, err := L.Match(*parms)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			resp := ErrorResponse(r.ID, 1, err.Error())
			if err := ToWriter(s.key, w, resp); err != nil {
				s.HTTPError(w, err)
			}
			return
		}
		rpcResp, err := NewResponse[map[string]map[string]int](r.ID, "loggers", &res)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		if err := ToWriter(s.key, w, rpcResp); err != nil {
			w.WriteHeader(http.StatusOK)
			resp := ErrorResponse(r.ID, 2, err.Error())
			if err := ToWriter(s.key, w, resp); err != nil {
				s.HTTPError(w, err)
			}
			return
		}
		return

	default:
		// -32601 is from jsonrpc 2.0
		resp := ErrorResponse(r.ID, -32601,
			fmt.Sprintf("%q: invalid method", r.Method))
		w.WriteHeader(http.StatusOK)
		if err := ToWriter(s.key, w, resp); err != nil {
			s.HTTPError(w, err)
		}
		// construct error
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
