package network

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ApplyOutboundControl(device string, latencyMs, lossPercent int) error {
	line, err := ShowOutboundControl(device)
	if err != nil {
		return err
	}
	exists := !strings.Contains(line, "qdisc noqueue 0")

	verb := "add"
	if exists {
		verb = "replace"
	}
	cmd := exec.Command(
		"tc",
		"qdisc",
		verb,
		"dev",
		device,
		"root",
		"netem",
		"delay",
		fmt.Sprintf("%dms", latencyMs),
		"loss",
		fmt.Sprintf("%d%%", lossPercent),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ShowOutboundControl(device string) (string, error) {
	cmd := exec.Command(
		"tc",
		"qdisc",
		"show",
		"dev",
		device,
		"root",
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
