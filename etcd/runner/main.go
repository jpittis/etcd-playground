package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
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
		log.Fatalf("usage: %s <exec-path> <name> <config-path>", os.Args[0])
	}
	execPath := os.Args[1]
	name := os.Args[2]
	configPath := os.Args[3]

	cfg, err := readConfigFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config config-path=%s, config=%+v", configPath, cfg)

	peerConfig := cfg.PeerConfig[name]
	for device, latencyMs := range peerConfig.OutboundLatencyMs {
		err := applyOutboundLatency(device, latencyMs)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Starting etcd exec-path=%s, name=%s", execPath, name)
	err = startEtcd(execPath, name)
	if err != nil {
		log.Fatalf("Etcd exited with error: %v", err)
	} else {
		log.Println("Etcd exited gracefully")
	}
}

func startEtcd(path, name string) error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)

	cmd := exec.Command(
		path,
		"--name", name,
		"--advertise-client-urls", fmt.Sprintf("http://%s:2379", name),
		"--listen-client-urls", "http://0.0.0.0:2379",
		"--initial-advertise-peer-urls", fmt.Sprintf("http://%s:2380", name),
		"--listen-peer-urls", "http://0.0.0.0:2380",
		"--initial-cluster-token", "etcd-cluster",
		"--initial-cluster", "etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380",
		"--initial-cluster-state", "new",
		"--enable-pprof",
		"--logger=zap",
		"--log-outputs=stderr",
	)
	// TODO(jpittis): Add log searching.
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		term := <-sigc
		log.Println("Terminating etcd")
		err := cmd.Process.Signal(term)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return cmd.Wait()
}

func applyOutboundLatency(device string, latencyMs int) error {
	log.Printf("Applying outbound latency device=%s, latency-ms=%d", device, latencyMs)
	cmd := exec.Command(
		"tc",
		"qdisc",
		"add",
		"dev",
		device,
		"root",
		"netem",
		"delay",
		fmt.Sprintf("%dms", latencyMs),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
