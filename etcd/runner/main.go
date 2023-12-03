package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <path> <name>", os.Args[0])
	}
	path := os.Args[1]
	name := os.Args[2]
	log.Printf("Starting etcd path=%s, name=%s", path, name)
	err := startEtcd(path, name)
	if err != nil {
		log.Fatalf("Exited with error: %v", err)
	} else {
		log.Println("Exited gracefully")
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
