package server

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
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
	mux.Handle("/log", http.HandlerFunc(s.log))
	return mux
}

func (s *Server) etcd(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "%t\n", s.proc.Enabled())
		return

	} else if r.Method != http.MethodPost {
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
	if r.Method == http.MethodGet {
		device := r.URL.Query().Get("dev")
		if device == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, "Must specify device (dev)\n")
			return
		}

		line, err := network.ShowOutboundControl(device)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v\n", err)
			return
		}

		fmt.Fprintf(w, "%s\n", line)
		return

	} else if r.Method != http.MethodPost {
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

func (s *Server) log(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Must GET\n")
		return
	}

	out, err := exec.Command(
		"/etcd/bin/tools/etcd-dump-logs",
		fmt.Sprintf("/etcd/%s.etcd", s.proc.Name()),
	).Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	fmt.Fprint(w, string(out))
}
