package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
)

type Firebase struct {
}

func (f Firebase) ConfigureHost() bool {
	fmt.Println("configured firebase host")
	return true
}

func (f Firebase) CoAddHostnfigureHost(alias string, definitionFilePath string, stateFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "firebase",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      stateFilePath,
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

func (f Firebase) ProviderTemplate() []byte {
	return []byte{}
}

func (f Firebase) ProviderConfigTemplate() []byte {
	return []byte{}
}
