package server

import (
	"fmt"
	"log"
	"net/http"
	"runner/pkg/network"
	"runner/pkg/process"
	"strconv"
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
		fmt.Fprint(w, "Must POST\n")
		return
	}

	enabled := r.URL.Query().Get("enabled")
	if enabled != "true" && enabled != "false" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "Must specify enabled true or false\n")
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
		fmt.Fprintf(w, "%v\n", err)
		return
	}
}

func (s *Server) network(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Must POST\n")
		return
	}

	device := r.URL.Query().Get("dev")
	if device == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "Must specify device (dev)\n")
		return
	}

	latency := r.URL.Query().Get("delay")
	var latencyMs int
	if latency != "" {
		var err error
		latencyMs, err = strconv.Atoi(latency)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v\n", err)
			return
		}
	}

	loss := r.URL.Query().Get("loss")
	var lossPercent int
	if loss != "" {
		var err error
		lossPercent, err = strconv.Atoi(loss)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v\n", err)
			return
		}
	}

	log.Printf("Applying outbound control dev=%s delay=%d loss=%d\n",
		device, latencyMs, lossPercent)
	err := network.ApplyOutboundControl(device, latencyMs, lossPercent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
}
