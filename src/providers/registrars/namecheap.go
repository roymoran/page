package providers

import (
	"fmt"

	"builtonpage.com/main/cliinit"
)

type Namecheap struct {
}

func (n Namecheap) ConfigureRegistrar() bool {
	fmt.Println("configured page registrar")
	// Does domain exist on registrar account? If not acquire.
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
