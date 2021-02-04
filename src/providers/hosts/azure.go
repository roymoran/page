package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
)

type Azure struct {
}

var AzureTerraformProvider = `
terraform {
	required_providers {
		azurerm = {
			source = "hashicorp/azurerm"
			version = "=2.44.0"
		}
	}
}
`

func (a Azure) ConfigureHost() bool {
	fmt.Println("configured azure host")
	return true
}

func (a Azure) AddHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "azure",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

func (a Azure) ProviderTemplate() []byte {
	return []byte{}
}

func (a Azure) ProviderConfigTemplate() []byte {
	return []byte{}
}
