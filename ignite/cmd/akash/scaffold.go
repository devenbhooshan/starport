package akash

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/cmd/akash/templates"
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
	log.Println("Starting scaffolding...")

	allFiles := []struct {
		fileName     string
		fileType     string
		fileContents string
	}{
		{
			fileName:     "dockerfileWeb.tpl",
			fileType:     "docker",
			fileContents: templates.FSAkashDockerFileWeb,
		},
		{
			fileName:     "dockerfileChain.tpl",
			fileType:     "docker",
			fileContents: templates.FSAkashDockerFileChain,
		},
		{
			fileName:     "deploy-chain.yml.tpl",
			fileType:     "SDL",
			fileContents: templates.FSAkashSDLChain,
		},
		{
			fileName:     "deploy-web.yml.tpl",
			fileType:     "SDL",
			fileContents: templates.FSAkashSDLWeb,
		},
	}

	for _, file := range allFiles {
		fileName := file.fileName
		fileType := file.fileType
		fileContents := file.fileContents
		err := createFileFromTemplate(fileName, fileContents, fileType)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
	}
	return nil
}

func createFileFromTemplate(fileName, fileContent, fileType string) error {
	tpl, err := template.New(fileName).Parse(fileContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	log.Println("Preparing template files...")

	folder := fmt.Sprintf("./akash/%s", fileType)
	out := path.Join(folder, strings.TrimSuffix(fileName, ".tpl"))

	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer outFile.Close()

	log.Printf("Creating %s file out of template file \n", fileType)

	err = tpl.ExecuteTemplate(outFile, fileName, nil)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
