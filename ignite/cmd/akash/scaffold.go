package akash

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewScaffold returns a command that scafolds the config for a network, default is AKASH
func NewScaffold() *cobra.Command {

	c := &cobra.Command{
		Use:   "scaffold testnet launch configs for a network, default is AKASH",
		Short: "scaffold a new account configs for a network, default is AKASH",
		RunE:  akashScaffoldHandler,
	}

	c.AddCommand()

	return c
}

func akashScaffoldHandler(cmd *cobra.Command, args []string) error {
	fmt.Println("scaffold a new account configs for a network, default is AKASH")
	return errors.New("not implemented")
}
