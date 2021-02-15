package providers

import (
	"encoding/json"
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Azure struct {
	HostName string
}

var azureProviderTemplate ProviderTemplate = ProviderTemplate{
	Terraform: RequiredProviders{
		RequiredProvider: map[string]Provider{
			"azurerm": {
				Source:  "hashicorp/azurerm",
				Version: "=2.44.0",
			},
		},
	},
}

func (a Azure) ConfigureAuth() error {
	return nil
}

func (a Azure) ConfigureHost(hostAlias string, registrarAlias string, templatePath string, page definition.PageDefinition) error {
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

// ProviderTemplate returns a byte slice that represents
// a template for creating an azure host
func (a Azure) ProviderTemplate() []byte {
	file, _ := json.MarshalIndent(azureProviderTemplate, "", " ")
	return file
}

// ProviderConfigTemplate returns a byte slice that represents
// configuration settings for the azure provider.
func (a Azure) ProviderConfigTemplate() []byte {
	return []byte{}
}
