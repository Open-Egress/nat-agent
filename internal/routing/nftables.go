package routing

import (
	"fmt"

	"github.com/google/nftables"
)

// NftablesConn defines the subset of *nftables.Conn methods used.
type NftablesConn interface {
	AddTable(t *nftables.Table) *nftables.Table
	FlushTable(t *nftables.Table)
	Flush() error
	CloseLasting() error
	AddChain(c *nftables.Chain) *nftables.Chain
}

// Ensure *nftables.Conn implements NftablesConn
var _ NftablesConn = (*nftables.Conn)(nil)

var NewNftablesConn = func() (NftablesConn, error) {
	return nftables.New()
}

func NftableApply() error {
	table := []*nftables.Table{
		{
			Family: nftables.TableFamilyINet,
			Name:   "nat",
		}, {
			Family: nftables.TableFamilyINet,
			Name:   "filter",
		},
	}

	return NftableCreateTables(table)
}

func NftableCreateTables(tables []*nftables.Table) error {
	conn, err := NewNftablesConn()
	if err != nil {
		return fmt.Errorf("failed to create nftables connection: %w", err)
	}

	for _, table := range tables {
		addedTable := conn.AddTable(table)
		conn.FlushTable(addedTable)
	}

	if err := conn.Flush(); err != nil {
		return fmt.Errorf("error while sending buffered commands to nftables: %w", err)
	}

	if err := conn.CloseLasting(); err != nil {
		return fmt.Errorf("error while closing nftables connection: %w", err)
	}

	return nil
}

func NftableCreateChain(tableName []string) error {
	conn, err := NewNftablesConn()
	if err != nil {
		return fmt.Errorf("failed to create nftables connection: %w", err)
	}

	for _, name := range tableName {
		conn.AddChain(&nftables.Chain{
			Name: "postrouting",
			Table: &nftables.Table{
				Name:   name,
				Family: nftables.TableFamilyIPv4,
			},
		})
	}

	if err := conn.Flush(); err != nil {
		return fmt.Errorf("error while sending buffered commands to nftables: %w", err)
	}

	return nil
}
