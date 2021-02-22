package hosts

import (
	"encoding/json"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Google struct {
	HostName string
}

var googleProviderTemplate ProviderTemplate = ProviderTemplate{
	Terraform: RequiredProviders{
		RequiredProvider: map[string]Provider{
			"google": {
				Source:  "hashicorp/google",
				Version: "3.56.0",
			},
		},
	},
}

func (g Google) ConfigureAuth() error {
	return nil
}

func (g Google) ConfigureHost(alias string, templatePath string, page definition.PageDefinition) error {
	return nil
}

func (g Google) CoAddHostnfigureHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "google",
		Credentials:      cliinit.Credentials{},
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// ProviderTemplate returns a byte slice that represents
// a template for creating an google host
func (g Google) ProviderTemplate() []byte {
	file, _ := json.MarshalIndent(azureProviderTemplate, "", " ")
	return file
}

// ProviderConfigTemplate returns a byte slice that represents
// configuration settings for the google provider.
func (g Google) ProviderConfigTemplate() []byte {
	return []byte{}
}
