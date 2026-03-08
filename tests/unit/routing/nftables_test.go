package routing_test

import (
	"errors"
	"testing"

	"nat-agent/internal/routing"

	"github.com/google/nftables"
)

type mockNftablesConn struct {
	addTableFunc     func(t *nftables.Table) *nftables.Table
	flushTableFunc   func(t *nftables.Table)
	flushFunc        func() error
	closeLastingFunc func() error
	addChainFunc     func(c *nftables.Chain) *nftables.Chain
}

func (m *mockNftablesConn) AddTable(t *nftables.Table) *nftables.Table { return m.addTableFunc(t) }
func (m *mockNftablesConn) FlushTable(t *nftables.Table)              { m.flushTableFunc(t) }
func (m *mockNftablesConn) Flush() error                             { return m.flushFunc() }
func (m *mockNftablesConn) CloseLasting() error                      { return m.closeLastingFunc() }
func (m *mockNftablesConn) AddChain(c *nftables.Chain) *nftables.Chain { return m.addChainFunc(c) }

func TestNftableCreateTables(t *testing.T) {
	oldNewConn := routing.NewNftablesConn
	defer func() { routing.NewNftablesConn = oldNewConn }()

	tests := []struct {
		name        string
		newConnErr  error
		flushErr    error
		closeErr    error
		expectedErr bool
	}{
		{
			name:        "success",
			newConnErr:  nil,
			flushErr:    nil,
			closeErr:    nil,
			expectedErr: false,
		},
		{
			name:        "new connection error",
			newConnErr:  errors.New("new connection error"),
			expectedErr: true,
		},
		{
			name:        "flush error",
			flushErr:    errors.New("flush error"),
			expectedErr: true,
		},
		{
			name:        "close error",
			closeErr:    errors.New("close error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routing.NewNftablesConn = func() (routing.NftablesConn, error) {
				if tt.newConnErr != nil {
					return nil, tt.newConnErr
				}
				return &mockNftablesConn{
					addTableFunc:     func(t *nftables.Table) *nftables.Table { return t },
					flushTableFunc:   func(t *nftables.Table) {},
					flushFunc:        func() error { return tt.flushErr },
					closeLastingFunc: func() error { return tt.closeErr },
				}, nil
			}

			err := routing.NftableCreateTables([]*nftables.Table{{Name: "test"}})
			if (err != nil) != tt.expectedErr {
				t.Errorf("NftableCreateTables() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestNftableCreateChain(t *testing.T) {
	oldNewConn := routing.NewNftablesConn
	defer func() { routing.NewNftablesConn = oldNewConn }()

	tests := []struct {
		name        string
		newConnErr  error
		flushErr    error
		expectedErr bool
	}{
		{
			name:        "success",
			newConnErr:  nil,
			flushErr:    nil,
			expectedErr: false,
		},
		{
			name:        "flush error",
			flushErr:    errors.New("flush error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routing.NewNftablesConn = func() (routing.NftablesConn, error) {
				if tt.newConnErr != nil {
					return nil, tt.newConnErr
				}
				return &mockNftablesConn{
					addChainFunc: func(c *nftables.Chain) *nftables.Chain { return c },
					flushFunc:    func() error { return tt.flushErr },
				}, nil
			}

			err := routing.NftableCreateChain([]*nftables.Chain{{Name: "test"}})
			if (err != nil) != tt.expectedErr {
				t.Errorf("NftableCreateChain() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
