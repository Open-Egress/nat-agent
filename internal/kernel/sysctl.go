package kernel

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// OSOps defines the interface for OS-level operations needed by IPForward.
type OSOps interface {
	ReadFile(name string) ([]byte, error)
	RunCommand(name string, arg ...string) error
}

// RealOS implements OSOps with real system calls.
type RealOS struct{}

func (RealOS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (RealOS) RunCommand(name string, arg ...string) error {
	return exec.Command(name, arg...).Run()
}

var Ops OSOps = RealOS{}

func IPForward() error {
	return IPForwardInternal(Ops)
}

func IPForwardInternal(ops OSOps) error {
	path := "/proc/sys/net/ipv4/ip_forward"
	desiredValue := "1"

	// Read current state
	content, err := ops.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	currentValue := strings.TrimSpace(string(content))
	slog.Info("Current kernel parameter state", "path", path, "value", currentValue)

	if currentValue == desiredValue {
		slog.Info("Parameter is already set correctly, skipping write", "path", path)
		return nil
	}

	slog.Warn("Parameter needs update, proceeding...", "path", path)

	// Update the parameter
	if err := ops.RunCommand("sysctl", "-w", "net.ipv4.ip_forward=1"); err != nil {
		return fmt.Errorf("failed to execute sysctl for %s: %w", path, err)
	}

	slog.Info("Kernel parameter updated successfully", "path", path)
	return nil
}
