package registrars

import (
	"pagecli.com/main/cliinit"
	"pagecli.com/main/definition"
	"pagecli.com/main/providers/hosts"
)

type Page struct {
}

func (p Page) ConfigureAuth() (cliinit.Credentials, error) {
	return cliinit.Credentials{}, nil
}
func (p Page) ConfigureDNS(registrarAlias string, hostAlias string, pageConfig definition.PageDefinition) error {
	return nil
}

func (p Page) ProviderDefinition() (string, hosts.Provider) {
	// TODO: Modify this
	return "page", hosts.Provider{Version: "1.7.0", Source: "page/page"}
}

func (p Page) AddRegistrar(alias string, credentials cliinit.Credentials) error {
	provider := cliinit.ProviderConfig{
		Type:             "registrar",
		Alias:            alias,
		Name:             "page",
		Credentials:      credentials,
		Default:          true,
		TfDefinitionPath: "",
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}
