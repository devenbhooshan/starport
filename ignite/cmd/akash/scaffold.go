package akash

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

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

	base := filepath.Join("starport/cmd/akash/templates", "docker")
	fileFullPath := filepath.Join(base, "dockerfileChain.tpl")

	template, err := template.ParseFiles(fileFullPath)
	if err != nil {
		return err
	}

	out := strings.TrimSuffix(filepath.Base(fileFullPath), ".tpl")

	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}
	defer f.Close()

	err = template.Execute(f, nil)
	if err != nil {
		return err
	}

	return errors.New("not implemented")
}
