package network

import (
	"fmt"
	"os"
	"os/exec"
)

func ApplyOutboundLatency(device string, latencyMs int) error {
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
