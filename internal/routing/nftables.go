package routing

import (
	"fmt"
	"log/slog"

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
	tables := []*nftables.Table{
		{
			Family: nftables.TableFamilyINet,
			Name:   "nat",
		}, {
			Family: nftables.TableFamilyINet,
			Name:   "filter",
		},
	}

	policyAccept := nftables.ChainPolicyAccept
	chains := []*nftables.Chain{
		{
			Name:     "postrouting",
			Table:    tables[0],
			Type:     nftables.ChainTypeNAT,
			Hooknum:  nftables.ChainHookPostrouting,
			Priority: nftables.ChainPriorityNATSource,
		},
		{
			Name:     "forward",
			Table:    tables[1],
			Type:     nftables.ChainTypeFilter,
			Hooknum:  nftables.ChainHookForward,
			Priority: nftables.ChainPriorityFilter,
			Policy:   &policyAccept,
		},
	}

	if err := NftableCreateTables(tables); err != nil {
		slog.Error("failed to create nftables tables", "error", err)
		return err
	}

	if err := NftableCreateChain(chains); err != nil {
		slog.Error("failed to create nftables chains", "error", err)
		return err
	}

	slog.Info("nftables configuration applied successfully")
	return nil
}

func NftableCreateTables(tables []*nftables.Table) error {
	conn, err := NewNftablesConn()
	if err != nil {
		return fmt.Errorf("failed to create nftables connection: %w", err)
	}

	for _, table := range tables {
		slog.Info("Adding nftables table", "name", table.Name, "family", table.Family)
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

func NftableCreateChain(chains []*nftables.Chain) error {
	conn, err := NewNftablesConn()
	if err != nil {
		return fmt.Errorf("failed to create nftables connection: %w", err)
	}

	for _, chain := range chains {
		tableName := "unknown"
		if chain.Table != nil {
			tableName = chain.Table.Name
		}
		slog.Info("Adding nftables chain", "name", chain.Name, "table", tableName)
		conn.AddChain(chain)
	}

	if err := conn.Flush(); err != nil {
		return fmt.Errorf("error while sending buffered commands to nftables: %w", err)
	}

	return nil
}
