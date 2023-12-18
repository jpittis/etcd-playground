package process

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var (
	ErrAlreadyStarted = errors.New("already started")
	ErrAlreadyStopped = errors.New("already stopped")
)

type Process struct {
	sync.Mutex
	execPath string
	nodeName string
	cmd      *exec.Cmd
}

func NewProcess(execPath, nodeName string) *Process {
	return &Process{
		execPath: execPath,
		nodeName: nodeName,
	}
}

func (p *Process) Start() error {
	p.Lock()
	defer p.Unlock()
	if p.cmd != nil {
		return ErrAlreadyStarted
	}

	cmd := buildCmd(p.execPath, p.nodeName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	p.cmd = cmd
	return nil
}

func (p *Process) Stop() error {
	p.Lock()
	defer p.Unlock()
	if p.cmd == nil {
		return ErrAlreadyStopped
	}

	err := p.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	err = p.cmd.Wait()
	p.cmd = nil
	return err
}

func buildCmd(execPath, nodeName string) *exec.Cmd {
	return exec.Command(
		execPath,
		"--name", nodeName,
		"--advertise-client-urls", fmt.Sprintf("http://%s:2379", nodeName),
		"--listen-client-urls", "http://0.0.0.0:2379",
		"--initial-advertise-peer-urls", fmt.Sprintf("http://%s:2380", nodeName),
		"--listen-peer-urls", "http://0.0.0.0:2380",
		"--initial-cluster-token", "etcd-cluster",
		"--initial-cluster",
		"etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380",
		"--initial-cluster-state", "new",
		"--enable-pprof",
		"--logger=zap",
		"--log-outputs=stderr",
	)
}
