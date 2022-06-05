package rpc

import (
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/scott-cotton/L"
)

type Server struct {
	key       []byte
	addr      string
	path      string
	log, warn L.Logger
}

func logConfig() *L.Config {
	cfg := L.NewConfig(".pkg", ".method")
	cfg.W = os.Stdout
	cfg.F = L.JSONFmter()
	cfg.E = L.EPanic
	cfg.Post = append(cfg.Post, L.Pkg())
	return cfg
}

func NewServer(key, addr, path string) *Server {
	srv := &Server{
		key:  []byte(key),
		addr: addr,
		path: path,
		log:  L.New(logConfig()),
	}
	srv.warn = srv.log.With(".warn", 1)
	return srv
}

func (s *Server) Serve() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	path := filepath.Join(s.path, "L")
	mux.HandleFunc(path, s.Handler)
	return http.Serve(ln, mux)
}
