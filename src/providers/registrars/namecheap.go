package providers

import "builtonpage.com/main/cliinit"

type Namecheap struct {
}

func (n Namecheap) RegisterDomain() bool {
	return true
}

func (n Namecheap) ConfigureDns() bool {
	return true
}

func (n Namecheap) ConfigureRegistrar(alias string) error {
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
