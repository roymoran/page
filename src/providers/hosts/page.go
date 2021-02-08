package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type PageHost struct {
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

func (p PageHost) ProviderTemplate() []byte {
	return []byte{}
}

func (p PageHost) ProviderConfigTemplate() []byte {
	return []byte{}
}
