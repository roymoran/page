package providers

import "builtonpage.com/main/cliinit"

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

func (a Azure) Deploy() bool {
	return true
}

func (a Azure) ConfigureHost(alias string, definitionFilePath string, stateFilePath string) error {
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
	return nil
}

func (a Azure) HostProviderDefinition() []byte {
	return []byte{}
}
