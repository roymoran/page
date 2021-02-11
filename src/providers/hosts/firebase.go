package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Firebase struct {
	HostName string
}

func (f Firebase) ConfigureAuth() error {
	return nil
}

func (f Firebase) ConfigureHost(alias string, templatePath string, page definition.PageDefinition) error {
	fmt.Println("configured firebase host")
	return nil
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

// ProviderInfo returns the provider info
// needed to download the terraform plugin
// for firebase
func (f Firebase) ProviderInfo() Provider {
	return Provider{
		Source:  "hashicorp/google",
		Version: "3.56.0",
	}
}

func (f Firebase) ProviderConfigTemplate() []byte {
	return []byte{}
}
