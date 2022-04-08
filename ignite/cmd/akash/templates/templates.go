package templates

import (
	_ "embed"
)

var (
	//go:embed docker/dockerfileWeb.tpl
	FSAkashDockerFileWeb string

	//go:embed docker/dockerfileChain.tpl
	FSAkashDockerFileChain string

	//go:embed SDL/deploy-web.yml.tpl
	FSAkashSDLWeb string

	//go:embed SDL/deploy-chain.yml.tpl
	FSAkashSDLChain string
)
