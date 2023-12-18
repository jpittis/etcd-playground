package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runner/pkg/network"
	"runner/pkg/process"
	"runner/pkg/server"
	"syscall"

	"gopkg.in/yaml.v3"
)

type config struct {
	PeerConfig map[string]peerConfig `yaml:"peer_config"`
}

type peerConfig struct {
	OutboundLatencyMs map[string]int `yaml:"outbound_latency_ms"`
}

func readConfigFromFile(path string) (config, error) {
	cfg := config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("usage: %s <exec-path> <node-name> <config-path>", os.Args[0])
	}
	execPath := os.Args[1]
	nodeName := os.Args[2]
	configPath := os.Args[3]

	cfg, err := readConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	log.Printf("Loaded config config-path=%s, config=%+v", configPath, cfg)

	peerConfig := cfg.PeerConfig[nodeName]
	for device, latencyMs := range peerConfig.OutboundLatencyMs {
		log.Printf("Applying outbound latency device=%s, latency-ms=%d", device, latencyMs)
		err := network.ApplyOutboundControl(device, latencyMs, 0)
		if err != nil {
			log.Fatalf("Failed to apply outbound latency: %v", err)
		}
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)

	log.Printf("Starting etcd exec-path=%s, node-name=%s", execPath, nodeName)
	proc := process.NewProcess(execPath, nodeName)
	err = proc.Start()
	if err != nil {
		log.Fatalf("Failed to start etcd: %v", err)
	}

	srv := server.NewServer(proc)
	mux := srv.NewServeMux()
	go func() {
		err := http.ListenAndServe("0.0.0.0:3333", mux)
		log.Fatalf("ListenAndServe exited: %v", err)
	}()

	done := make(chan struct{}, 0)
	go func() {
		<-sigc
		log.Println("Terminating etcd")
		err := proc.Stop()
		if err != nil {
			log.Fatalf("Failed to stop etcd: %v", err)
		}
		close(done)
	}()
	<-done
}
