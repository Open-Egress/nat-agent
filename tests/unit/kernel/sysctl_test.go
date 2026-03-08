package kernel_test

import (
	"errors"
	"testing"

	"nat-agent/internal/kernel"
)

type mockOS struct {
	readFileFunc   func(name string) ([]byte, error)
	runCommandFunc func(name string, arg ...string) error
}

func (m *mockOS) ReadFile(name string) ([]byte, error) {
	return m.readFileFunc(name)
}

func (m *mockOS) RunCommand(name string, arg ...string) error {
	return m.runCommandFunc(name, arg...)
}

func TestIPForward(t *testing.T) {
	tests := []struct {
		name           string
		readFileRes    []byte
		readFileErr    error
		runCommandErr  error
		expectedErr    bool
		commandVisited bool
	}{
		{
			name:           "already enabled",
			readFileRes:    []byte("1\n"),
			readFileErr:    nil,
			runCommandErr:  nil,
			expectedErr:    false,
			commandVisited: false,
		},
		{
			name:           "needs enabling",
			readFileRes:    []byte("0\n"),
			readFileErr:    nil,
			runCommandErr:  nil,
			expectedErr:    false,
			commandVisited: true,
		},
		{
			name:           "read error",
			readFileRes:    nil,
			readFileErr:    errors.New("read error"),
			runCommandErr:  nil,
			expectedErr:    true,
			commandVisited: false,
		},
		{
			name:           "command error",
			readFileRes:    []byte("0\n"),
			readFileErr:    nil,
			runCommandErr:  errors.New("command error"),
			expectedErr:    true,
			commandVisited: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := false
			mock := &mockOS{
				readFileFunc: func(name string) ([]byte, error) {
					return tt.readFileRes, tt.readFileErr
				},
				runCommandFunc: func(name string, arg ...string) error {
					visited = true
					return tt.runCommandErr
				},
			}

			err := kernel.IPForwardInternal(mock)
			if (err != nil) != tt.expectedErr {
				t.Errorf("IPForward() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if visited != tt.commandVisited {
				t.Errorf("IPForward() command visited = %v, expected %v", visited, tt.commandVisited)
			}
		})
	}
}
