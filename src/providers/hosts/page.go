package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
)

type PageHost struct {
}

func (p PageHost) ConfigureHost(alias string) error {
	fmt.Println("configured page host")
	return nil
}

func (p PageHost) AddHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "page",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
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
