package server

import (
	"net/http"
	"runner/pkg/process"
)

type Server struct {
	proc *process.Process
}

func NewServer(proc *process.Process) *Server {
	return &Server{
		proc: proc,
	}
}

func (s *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/etcd", http.HandlerFunc(s.etcd))
	mux.Handle("/network", http.HandlerFunc(s.network))
	return mux
}

func (s *Server) etcd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	enabled := r.URL.Query().Get("enabled")
	if enabled != "true" && enabled != "false" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var err error
	if enabled == "true" {
		err = s.proc.Start()
	} else {
		err = s.proc.Stop()
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) network(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
