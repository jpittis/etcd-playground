package server

import (
	"fmt"
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
		fmt.Fprint(w, "Must POST")
		return
	}

	enabled := r.URL.Query().Get("enabled")
	if enabled != "true" && enabled != "false" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "Must specify enabled true or false")
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
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func (s *Server) network(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Must POST")
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Unimplemented")
}
