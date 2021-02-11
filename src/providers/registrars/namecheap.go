package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
	"builtonpage.com/main/definition"
)

type Namecheap struct {
}

func (n Namecheap) ConfigureAuth() error {
	fmt.Println("configured namecheap registrar auth")
	return nil
}

func (n Namecheap) ConfigureRegistrar(pageConfig definition.PageDefinition) bool {
	fmt.Println("configured namecheap registrar")
	// TODO: Generate SSL cert and validate against domain
	// with DNS validation
	return true
}

func (n Namecheap) ConfigureDns() bool {
	return true
}

func (n Namecheap) AddRegistrar(alias string) error {
	provider := cliinit.ProviderConfig{
		Type:             "registrar",
		Alias:            alias,
		Name:             "namecheap",
		Auth:             "tbd",
		Default:          true,
		TfDefinitionPath: "",
		TfStatePath:      "",
	}

	addProviderErr := cliinit.AddProvider(provider)
	return addProviderErr
}
