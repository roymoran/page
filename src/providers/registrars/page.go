package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Page struct {
}

func (p Page) ConfigureAuth() error {
	fmt.Println("configured namecheap registrar auth")
	return nil
}
func (p Page) ConfigureRegistrar(pageConfig definition.PageDefinition) bool {
	fmt.Println("configured page registrar")
	return true
}

func (p Page) ConfigureDns() bool {
	return true
}

func (p Page) AddRegistrar(alias string) error {
	provider := cliinit.ProviderConfig{
		Type:             "registrar",
		Alias:            alias,
		Name:             "page",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: "",
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}
