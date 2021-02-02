package providers

import "builtonpage.com/main/cliinit"

type Page struct {
}

func (p Page) RegisterDomain() bool {
	return true
}

func (p Page) ConfigureDns() bool {
	return true
}

func (p Page) ConfigureRegistrar(alias string) error {
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
