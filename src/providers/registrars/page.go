package registrars

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
	"builtonpage.com/main/providers/hosts"
)

type Page struct {
}

func (p Page) ConfigureAuth() (cliinit.Credentials, error) {
	fmt.Println("configured page registrar auth")
	return cliinit.Credentials{}, nil
}
func (p Page) ConfigureRegistrar(registrarAlias string, hostAlias string, pageConfig definition.PageDefinition) error {
	fmt.Println("configured page registrar")
	return nil
}

func (p Page) ProviderDefinition() (string, hosts.Provider) {
	// TODO: Modify this
	return "page", hosts.Provider{Version: "1.7.0", Source: "robgmills/namecheap"}
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
