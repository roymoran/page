package hosts

import (
	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type PageHost struct {
	HostName string
}

func (p PageHost) ConfigureAuth() error {
	return nil
}

func (p PageHost) ConfigureHost(hostAlias string, templatePath string, page definition.PageDefinition) error {
	return nil
}

func (p PageHost) AddHost(alias string, definitionFilePath string) error {
	provider := cliinit.ProviderConfig{
		Type:             "host",
		Alias:            alias,
		Name:             "page",
		Credentials:      cliinit.Credentials{},
		Default:          true,
		TfDefinitionPath: definitionFilePath,
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}

// ProviderConfigInfo returns a mapping representing configuration
// settings for the page terraform provider
func (p PageHost) ProviderConfigInfo() interface{} {
	return map[string]interface{}{
		p.HostName: map[string]interface{}{
			"region":     "us-east-1",
			"access_key": accessKey,
			"secret_key": secretKey,
		},
	}
}

// ProviderInfo returns a byte slice that represents
// a template for creating an aws host
func (p PageHost) ProviderInfo() Provider {
	return Provider{
		Source:  "hashicorp/aws",
		Version: "3.25.0",
	}
}

// ProviderTemplate returns a byte slice that represents
// a template for creating an google host
func (p PageHost) ProviderTemplate() []byte {
	return []byte{}
}

// ProviderConfigTemplate returns a byte slice that represents
// configuration settings for the google provider.
func (p PageHost) ProviderConfigTemplate() []byte {
	return []byte{}
}
