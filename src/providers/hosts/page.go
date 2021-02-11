package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type PageHost struct {
	HostName string
}

func (p PageHost) ConfigureAuth() error {
	return nil
}

func (p PageHost) ConfigureHost(alias string, templatePath string, page definition.PageDefinition) error {
	fmt.Println("configured page host")
	return nil
}

func (p PageHost) AddHost(alias string, definitionFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "page",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// ProviderInfo returns a byte slice that represents
// a template for creating an aws host
func (p PageHost) ProviderInfo() Provider {
	return Provider{
		Source:  "hashicorp/aws",
		Version: "3.25.0",
	}
}

func (p PageHost) ProviderConfigTemplate() []byte {
	return []byte{}
}
