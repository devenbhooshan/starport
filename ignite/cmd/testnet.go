package ignitecmd

import (
	"github.com/ignite-hq/cli/ignite/cmd/akash"
	"github.com/spf13/cobra"
)

// NewTestnet returns a command that groups sub commands related to generating and launching on testnet.
func NewTestnet() *cobra.Command {
	c := &cobra.Command{
		Use:     "testnet [command]",
		Short:   "launch chain on testnet, default network AKASH",
		Long:    "launch chain on testnet, default network AKASH",
		Aliases: []string{"t"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(
		akash.NewScaffold(),
		akash.NewLaunch(),
	)

	return c
}
