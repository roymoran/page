package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Azure struct {
	HostName string
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

func (a Azure) ConfigureAuth() error {
	return nil
}

func (a Azure) ConfigureHost(alias string, templatePath string, page definition.PageDefinition) error {
	fmt.Println("configured azure host")
	return nil
}

func (a Azure) AddHost(alias string, definitionFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "azure",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// ProviderInfo returns the provider info
// needed to download the terraform plugin
// for azure
func (a Azure) ProviderInfo() Provider {
	return Provider{
		Source:  "hashicorp/azurerm",
		Version: "=2.47.0",
	}
}

func (a Azure) ProviderConfigTemplate() []byte {
	return []byte{}
}
